import customtkinter as ctk
import tkinter as tk
from tkinter import messagebox
import json
import os
import sys

# Set theme and color options
ctk.set_appearance_mode("System")
ctk.set_default_color_theme("blue")

# Scaling for high-DPI (4K) monitors. 
ctk.set_widget_scaling(2.0)
ctk.set_window_scaling(2.0)

CONFIG_PATH = os.path.expanduser("~/.factory/config.json")

class ConfigEditorApp(ctk.CTk):
    def __init__(self):
        super().__init__()

        self.title("Droid Config Editor")
        self.geometry("1400x900")

        self.config_data = {"custom_models": []}
        self.current_index = None
        self.selected_indices = set()  # For multi-select deletion

        # Main Layout Configuration
        self.grid_columnconfigure(1, weight=1)
        self.grid_rowconfigure(0, weight=1)

        # --- Left Panel (Sidebar) ---
        self.sidebar_frame = ctk.CTkFrame(self, width=300, corner_radius=0)
        self.sidebar_frame.grid(row=0, column=0, sticky="nsew")
        self.sidebar_frame.grid_propagate(False) 

        # 1. Logo (Top)
        self.logo_label = ctk.CTkLabel(self.sidebar_frame, text="Droid Config", font=ctk.CTkFont(size=24, weight="bold"))
        self.logo_label.pack(pady=(20, 10), padx=20)

        # 2. New Model Button (Top, below logo)
        self.btn_new_model = ctk.CTkButton(self.sidebar_frame, text="New Model", command=self.new_model)
        self.btn_new_model.pack(pady=10, padx=20, fill="x")

        # 3. Select All Checkbox
        self.var_select_all = tk.BooleanVar(value=False)
        self.chk_select_all = ctk.CTkCheckBox(self.sidebar_frame, text="Select All", 
                                               variable=self.var_select_all, command=self.toggle_select_all)
        self.chk_select_all.pack(pady=(10, 0), padx=20, anchor="w")

        # 4. Delete Selected Button
        self.btn_delete_selected = ctk.CTkButton(self.sidebar_frame, text="Delete Selected", 
                                                  fg_color="red", hover_color="darkred",
                                                  command=self.delete_selected_models, state="disabled")
        self.btn_delete_selected.pack(pady=10, padx=20, fill="x")

        # 5. Scrollable List (Fills remaining space)
        self.scrollable_list = ctk.CTkScrollableFrame(self.sidebar_frame, label_text="Your Models")
        self.scrollable_list.pack(fill="both", expand=True, padx=20, pady=10)
        
        self.list_buttons = []
        self.list_checkboxes = []
        self.checkbox_vars = []

        # --- Right Panel (Main Content) ---
        self.main_frame = ctk.CTkFrame(self, corner_radius=20, fg_color="transparent")
        self.main_frame.grid(row=0, column=1, sticky="nsew", padx=20, pady=20)
        self.main_frame.grid_columnconfigure(0, weight=1)

        # Header / Status
        self.lbl_header = ctk.CTkLabel(self.main_frame, text="Edit Model", font=ctk.CTkFont(size=28, weight="bold"))
        self.lbl_header.grid(row=0, column=0, sticky="w", pady=(0, 20))
        
        # Status Label (Hidden by default, used for auto-save feedback)
        self.lbl_status = ctk.CTkLabel(self.main_frame, text="", font=ctk.CTkFont(size=14, slant="italic"), text_color="green")
        self.lbl_status.grid(row=0, column=1, sticky="e", pady=(0, 20), padx=20)

        # Form Container
        self.form_frame = ctk.CTkFrame(self.main_frame, corner_radius=15)
        self.form_frame.grid(row=1, column=0, columnspan=2, sticky="nsew")
        self.form_frame.grid_columnconfigure(1, weight=1)

        # Form Variables
        self.var_display_name = tk.StringVar()
        self.var_model = tk.StringVar()
        self.var_base_url = tk.StringVar()
        self.var_api_key = tk.StringVar()
        self.var_provider = tk.StringVar(value="anthropic")
        self.var_max_tokens = tk.StringVar(value="8192")

        # Create Form Fields
        self.create_form_entry(0, "Display Name", self.var_display_name, "Human-friendly name")
        self.create_form_entry(1, "Model ID", self.var_model, "e.g., gpt-4-turbo")
        self.create_form_entry(2, "Base URL", self.var_base_url, "API Endpoint")
        self.create_form_entry(3, "API Key", self.var_api_key, "Secret Key")
        
        # Provider Combobox
        lbl_prov = ctk.CTkLabel(self.form_frame, text="Provider:", font=ctk.CTkFont(size=16, weight="bold"))
        lbl_prov.grid(row=4, column=0, padx=20, pady=(15, 0), sticky="w")
        
        self.combo_provider = ctk.CTkComboBox(self.form_frame, variable=self.var_provider, 
                                              values=["anthropic", "openai", "generic-chat-completion-api"])
        self.combo_provider.grid(row=4, column=1, padx=20, pady=(15, 0), sticky="ew")

        # Max Tokens
        self.create_form_entry(5, "Max Tokens", self.var_max_tokens, "Max output length")

        # Action Buttons Area
        self.actions_frame = ctk.CTkFrame(self.main_frame, fg_color="transparent")
        self.actions_frame.grid(row=2, column=0, columnspan=2, sticky="ew", pady=20)

        self.btn_delete = ctk.CTkButton(self.actions_frame, text="Delete Model", fg_color="red", hover_color="darkred",
                                        command=self.delete_model, state="disabled")
        self.btn_delete.pack(side="left")

        self.btn_move_up = ctk.CTkButton(self.actions_frame, text="↑ Move Up", width=100,
                                         command=self.move_model_up, state="disabled")
        self.btn_move_up.pack(side="left", padx=(10, 5))

        self.btn_move_down = ctk.CTkButton(self.actions_frame, text="↓ Move Down", width=100,
                                           command=self.move_model_down, state="disabled")
        self.btn_move_down.pack(side="left", padx=5)

        self.btn_apply = ctk.CTkButton(self.actions_frame, text="Save Changes", command=self.apply_changes)
        self.btn_apply.pack(side="right")
        
        self.load_config()

    def create_form_entry(self, row, label_text, variable, placeholder=""):
        lbl = ctk.CTkLabel(self.form_frame, text=f"{label_text}:", font=ctk.CTkFont(size=16, weight="bold"))
        lbl.grid(row=row, column=0, padx=20, pady=(15, 0), sticky="w")
        
        entry = ctk.CTkEntry(self.form_frame, textvariable=variable, placeholder_text=placeholder)
        entry.grid(row=row, column=1, padx=20, pady=(15, 0), sticky="ew")

    def load_config(self):
        if not os.path.exists(CONFIG_PATH):
            try:
                os.makedirs(os.path.dirname(CONFIG_PATH), exist_ok=True)
            except OSError:
                pass
            self.config_data = {"custom_models": []}
        else:
            try:
                with open(CONFIG_PATH, 'r') as f:
                    self.config_data = json.load(f)
            except Exception as e:
                messagebox.showerror("Error", f"Failed to load config: {e}")
                self.config_data = {"custom_models": []}
        
        self.refresh_list()

    def refresh_list(self):
        for btn in self.list_buttons:
            btn.destroy()
        for chk in self.list_checkboxes:
            chk.destroy()
        self.list_buttons = []
        self.list_checkboxes = []
        self.checkbox_vars = []
        self.selected_indices.clear()
        self.var_select_all.set(False)
        self._update_delete_selected_button()

        models = self.config_data.get("custom_models", [])
        for idx, model in enumerate(models):
            display_name = model.get("model_display_name", "Unnamed")
            
            # Frame to hold checkbox and button
            item_frame = ctk.CTkFrame(self.scrollable_list, fg_color="transparent")
            item_frame.pack(fill="x", pady=2)
            
            # Checkbox for multi-select
            var = tk.BooleanVar(value=False)
            chk = ctk.CTkCheckBox(item_frame, text="", variable=var, width=24,
                                  command=lambda i=idx: self.toggle_model_selection(i))
            chk.pack(side="left", padx=(0, 5))
            self.checkbox_vars.append(var)
            self.list_checkboxes.append(item_frame)  # Store frame for destruction
            
            # Button for selection/editing
            btn = ctk.CTkButton(item_frame, text=display_name, fg_color="transparent", 
                                text_color=("gray10", "gray90"), hover_color=("gray70", "gray30"), anchor="w",
                                command=lambda i=idx: self.select_model(i))
            btn.pack(side="left", fill="x", expand=True)
            self.list_buttons.append(btn)
        
        if self.current_index is not None and self.current_index < len(models):
            self.select_model(self.current_index)
        else:
            self.new_model()

    def select_model(self, index):
        self.current_index = index
        model_data = self.config_data["custom_models"][index]
        
        self.var_display_name.set(model_data.get("model_display_name", ""))
        self.var_model.set(model_data.get("model", ""))
        self.var_base_url.set(model_data.get("base_url", ""))
        self.var_api_key.set(model_data.get("api_key", ""))
        self.var_provider.set(model_data.get("provider", "anthropic"))
        self.var_max_tokens.set(str(model_data.get("max_tokens", 8192)))

        self.btn_delete.configure(state="normal")
        self._update_move_buttons()
        self.lbl_header.configure(text=f"Editing: {self.var_display_name.get()}")
        self.lbl_status.configure(text="") # Clear status

        for i, btn in enumerate(self.list_buttons):
            if i == index:
                btn.configure(fg_color=("gray75", "gray25"))
            else:
                btn.configure(fg_color="transparent")

    def new_model(self):
        self.current_index = None
        self.var_display_name.set("")
        self.var_model.set("")
        self.var_base_url.set("")
        self.var_api_key.set("")
        self.var_provider.set("anthropic")
        self.var_max_tokens.set("8192")
        
        self.btn_delete.configure(state="disabled")
        self._update_move_buttons()
        self.lbl_header.configure(text="New Model")
        self.lbl_status.configure(text="")

        for btn in self.list_buttons:
            btn.configure(fg_color="transparent")

    def apply_changes(self):
        if not self.var_display_name.get().strip():
            messagebox.showwarning("Validation", "Display Name is required.")
            return
        
        try:
            max_tok = int(self.var_max_tokens.get())
        except ValueError:
            messagebox.showwarning("Validation", "Max Tokens must be an integer.")
            return

        model_entry = {
            "model_display_name": self.var_display_name.get(),
            "model": self.var_model.get(),
            "base_url": self.var_base_url.get(),
            "api_key": self.var_api_key.get(),
            "provider": self.var_provider.get(),
            "max_tokens": max_tok
        }

        if self.current_index is not None:
            self.config_data["custom_models"][self.current_index] = model_entry
        else:
            if "custom_models" not in self.config_data:
                self.config_data["custom_models"] = []
            self.config_data["custom_models"].append(model_entry)
            self.current_index = len(self.config_data["custom_models"]) - 1
        
        self._save_config() # Auto-save
        self.refresh_list()
        
        # Show status feedback
        self.lbl_status.configure(text="Changes Saved!", text_color="green")
        self.after(2000, lambda: self.lbl_status.configure(text=""))

    def delete_model(self):
        if self.current_index is not None:
            if messagebox.askyesno("Confirm", "Delete this model?"):
                del self.config_data["custom_models"][self.current_index]
                self.current_index = None
                self._save_config() # Auto-save
                self.refresh_list()
                self.lbl_status.configure(text="Model Deleted", text_color="red")
                self.after(2000, lambda: self.lbl_status.configure(text=""))

    def toggle_model_selection(self, index):
        """Toggle selection state for a model checkbox."""
        if self.checkbox_vars[index].get():
            self.selected_indices.add(index)
        else:
            self.selected_indices.discard(index)
        self._update_delete_selected_button()
        self._update_select_all_state()

    def toggle_select_all(self):
        """Select or deselect all models."""
        select_all = self.var_select_all.get()
        for idx, var in enumerate(self.checkbox_vars):
            var.set(select_all)
            if select_all:
                self.selected_indices.add(idx)
            else:
                self.selected_indices.discard(idx)
        self._update_delete_selected_button()

    def _update_select_all_state(self):
        """Update select all checkbox based on individual selections."""
        if not self.checkbox_vars:
            self.var_select_all.set(False)
        elif len(self.selected_indices) == len(self.checkbox_vars):
            self.var_select_all.set(True)
        else:
            self.var_select_all.set(False)

    def _update_delete_selected_button(self):
        """Enable/disable delete selected button based on selections."""
        if self.selected_indices:
            self.btn_delete_selected.configure(state="normal")
        else:
            self.btn_delete_selected.configure(state="disabled")

    def delete_selected_models(self):
        """Delete all selected models."""
        if not self.selected_indices:
            return
        
        count = len(self.selected_indices)
        if messagebox.askyesno("Confirm", f"Delete {count} selected model(s)?"):
            # Delete in reverse order to maintain correct indices
            for idx in sorted(self.selected_indices, reverse=True):
                del self.config_data["custom_models"][idx]
            
            self.current_index = None
            self.selected_indices.clear()
            self._save_config()
            self.refresh_list()
            self.lbl_status.configure(text=f"{count} Model(s) Deleted", text_color="red")
            self.after(2000, lambda: self.lbl_status.configure(text=""))

    def move_model_up(self):
        if self.current_index is not None and self.current_index > 0:
            models = self.config_data["custom_models"]
            models[self.current_index], models[self.current_index - 1] = \
                models[self.current_index - 1], models[self.current_index]
            self.current_index -= 1
            self._save_config()
            self.refresh_list()
            self.lbl_status.configure(text="Model Moved Up", text_color="green")
            self.after(2000, lambda: self.lbl_status.configure(text=""))

    def move_model_down(self):
        if self.current_index is not None:
            models = self.config_data["custom_models"]
            if self.current_index < len(models) - 1:
                models[self.current_index], models[self.current_index + 1] = \
                    models[self.current_index + 1], models[self.current_index]
                self.current_index += 1
                self._save_config()
                self.refresh_list()
                self.lbl_status.configure(text="Model Moved Down", text_color="green")
                self.after(2000, lambda: self.lbl_status.configure(text=""))

    def _update_move_buttons(self):
        if self.current_index is None:
            self.btn_move_up.configure(state="disabled")
            self.btn_move_down.configure(state="disabled")
        else:
            models = self.config_data.get("custom_models", [])
            # Enable/disable up button
            if self.current_index > 0:
                self.btn_move_up.configure(state="normal")
            else:
                self.btn_move_up.configure(state="disabled")
            # Enable/disable down button
            if self.current_index < len(models) - 1:
                self.btn_move_down.configure(state="normal")
            else:
                self.btn_move_down.configure(state="disabled")

    def _save_config(self):
        try:
            with open(CONFIG_PATH, 'w') as f:
                json.dump(self.config_data, f, indent=2)
        except Exception as e:
            messagebox.showerror("Error", f"Failed to save: {e}")

def main():
    app = ConfigEditorApp()
    app.mainloop()

if __name__ == "__main__":
    main()
