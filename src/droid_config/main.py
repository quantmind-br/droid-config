import tkinter as tk
from tkinter import messagebox, filedialog
import ttkbootstrap as ttk
from ttkbootstrap.constants import *
import json
import os
import sys

CONFIG_PATH = os.path.expanduser("~/.factory/config.json")

class ConfigEditorApp:
    def __init__(self, root):
        self.root = root
        self.root.title("Droid Config Editor")
        self.root.geometry("1024x768")
        
        # We let ttkbootstrap handle the basic styles (fonts, colors).
        # We can define custom styles if needed here, but usually the theme is enough.

        self.config_data = {"custom_models": []}
        self.current_index = None

        # Create main layout
        # Use a container for padding
        self.main_container = ttk.Frame(root, padding=20)
        self.main_container.pack(fill=BOTH, expand=True)

        self.paned_window = ttk.Panedwindow(self.main_container, orient=HORIZONTAL)
        self.paned_window.pack(fill=BOTH, expand=True)

        # Left panel: List of models
        self.left_frame = ttk.Frame(self.paned_window, width=300)
        self.paned_window.add(self.left_frame, weight=1)

        ttk.Label(self.left_frame, text="Models", font=("Helvetica", 16, "bold"), bootstyle="primary").pack(anchor=W, pady=(0, 15))
        
        # Standard Listbox doesn't fully support ttk styles, but we can set its fonts/colors to match loosely
        # or just keep it simple.
        self.model_listbox = tk.Listbox(self.left_frame, selectmode=tk.SINGLE, font=("Helvetica", 12), 
                                        exportselection=False, borderwidth=1, relief="flat", highlightthickness=1)
        self.model_listbox.pack(fill=BOTH, expand=True, padx=(0, 10))
        self.model_listbox.bind('<<ListboxSelect>>', self.on_select)

        # Action buttons on the left
        btn_container_left = ttk.Frame(self.left_frame)
        btn_container_left.pack(fill=X, pady=10, padx=(0, 10))

        self.btn_add = ttk.Button(btn_container_left, text="New Model", command=self.new_model, bootstyle="success-outline")
        self.btn_add.pack(fill=X, pady=(0, 5))
        
        self.btn_save_file = ttk.Button(btn_container_left, text="💾 Save to Disk", command=self.save_to_file, bootstyle="primary")
        self.btn_save_file.pack(fill=X)

        # Right panel: Edit form
        self.right_frame = ttk.Frame(self.paned_window, padding=(20, 0, 0, 0))
        self.paned_window.add(self.right_frame, weight=3)

        self.create_form()
        self.load_config()

    def create_form(self):
        # Variables
        self.var_display_name = tk.StringVar()
        self.var_model = tk.StringVar()
        self.var_base_url = tk.StringVar()
        self.var_api_key = tk.StringVar()
        self.var_provider = tk.StringVar()
        self.var_max_tokens = tk.IntVar(value=8192)

        # Form Container
        # Just use right_frame directly or a wrapper
        form_container = self.right_frame

        self.lbl_status = ttk.Label(form_container, text="Creating New Model", font=("Helvetica", 12, "italic"), bootstyle="info")
        self.lbl_status.grid(row=0, column=0, columnspan=3, pady=(0, 20), sticky=W)

        # Fields
        self.add_field(form_container, "Display Name", self.var_display_name, 1, 
                       "Human-friendly name shown in model selector")
        self.add_field(form_container, "Model ID", self.var_model, 2,
                       "Model identifier sent via API (e.g., gpt-5-codex)")
        self.add_field(form_container, "Base URL", self.var_base_url, 3,
                       "API endpoint base URL")
        self.add_field(form_container, "API Key", self.var_api_key, 4,
                       "Your API key for the provider")
        
        # Provider Dropdown
        ttk.Label(form_container, text="Provider", bootstyle="secondary").grid(row=5, column=0, sticky=W, pady=10)
        self.combo_provider = ttk.Combobox(form_container, textvariable=self.var_provider, state="readonly", font=("Helvetica", 11))
        self.combo_provider['values'] = ('anthropic', 'openai', 'generic-chat-completion-api')
        self.combo_provider.grid(row=5, column=1, sticky=EW, pady=10, ipady=4)
        ttk.Label(form_container, text="One of: anthropic, openai, or generic-chat-completion-api", 
                  font=("Arial", 9), bootstyle="secondary").grid(row=5, column=2, sticky=W, padx=10)

        self.add_field(form_container, "Max Tokens", self.var_max_tokens, 6,
                       "Maximum output tokens for model responses")

        # Form Actions
        btn_frame = ttk.Frame(form_container)
        btn_frame.grid(row=7, column=0, columnspan=3, pady=40, sticky=E)

        self.btn_delete = ttk.Button(btn_frame, text="Delete", command=self.delete_model, state=DISABLED, bootstyle="danger")
        self.btn_delete.pack(side=LEFT, padx=10)

        self.btn_apply = ttk.Button(btn_frame, text="Apply Changes", command=self.apply_changes, bootstyle="success")
        self.btn_apply.pack(side=LEFT, padx=10)

        # Configure grid
        form_container.columnconfigure(1, weight=1)

    def add_field(self, parent, label, variable, row, help_text):
        ttk.Label(parent, text=label, bootstyle="secondary").grid(row=row, column=0, sticky=W, pady=10)
        entry = ttk.Entry(parent, textvariable=variable, font=("Helvetica", 11))
        entry.grid(row=row, column=1, sticky=EW, pady=10, ipady=4)
        ttk.Label(parent, text=help_text, font=("Arial", 9), bootstyle="secondary").grid(row=row, column=2, sticky=W, padx=10)

    def load_config(self):
        if not os.path.exists(CONFIG_PATH):
            try:
                os.makedirs(os.path.dirname(CONFIG_PATH), exist_ok=True)
            except OSError:
                pass # Directory might exist
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
        self.model_listbox.delete(0, tk.END)
        for model in self.config_data.get("custom_models", []):
            self.model_listbox.insert(tk.END, model.get("model_display_name", "Unnamed"))
        
        if self.current_index is not None and self.current_index < len(self.config_data["custom_models"]):
            self.model_listbox.selection_set(self.current_index)
            self.load_model_to_form(self.config_data["custom_models"][self.current_index])
        else:
            self.new_model()

    def on_select(self, event):
        selection = self.model_listbox.curselection()
        if selection:
            self.current_index = selection[0]
            self.load_model_to_form(self.config_data["custom_models"][self.current_index])
            self.btn_delete.config(state=NORMAL)
            self.lbl_status.config(text=f"Editing: {self.var_display_name.get()}", bootstyle="success")

    def load_model_to_form(self, model_data):
        self.var_display_name.set(model_data.get("model_display_name", ""))
        self.var_model.set(model_data.get("model", ""))
        self.var_base_url.set(model_data.get("base_url", ""))
        self.var_api_key.set(model_data.get("api_key", ""))
        self.var_provider.set(model_data.get("provider", ""))
        self.var_max_tokens.set(model_data.get("max_tokens", 8192))
        self.lbl_status.config(text=f"Editing: {self.var_display_name.get()}", bootstyle="success")

    def new_model(self):
        self.current_index = None
        self.model_listbox.selection_clear(0, tk.END)
        self.var_display_name.set("")
        self.var_model.set("")
        self.var_base_url.set("")
        self.var_api_key.set("")
        self.var_provider.set("")
        self.var_max_tokens.set(8192)
        self.btn_delete.config(state=DISABLED)
        self.lbl_status.config(text="Creating New Model", bootstyle="info")

    def delete_model(self):
        if self.current_index is not None:
            confirm = messagebox.askyesno("Confirm", "Are you sure you want to delete this model?")
            if confirm:
                del self.config_data["custom_models"][self.current_index]
                self.current_index = None
                self.refresh_list()

    def apply_changes(self):
        # Validation
        if not self.var_display_name.get():
            messagebox.showwarning("Validation", "Display Name is required.")
            return
        if not self.var_model.get():
            messagebox.showwarning("Validation", "Model ID is required.")
            return
        if not self.var_provider.get():
            messagebox.showwarning("Validation", "Provider is required.")
            return
        
        try:
            max_tok = self.var_max_tokens.get()
        except:
            messagebox.showwarning("Validation", "Max Tokens must be a number.")
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
            # Update existing
            self.config_data["custom_models"][self.current_index] = model_entry
            status_msg = "Updated existing model."
        else:
            # Add new
            if "custom_models" not in self.config_data:
                self.config_data["custom_models"] = []
            self.config_data["custom_models"].append(model_entry)
            self.current_index = len(self.config_data["custom_models"]) - 1
            status_msg = "Added new model."
        
        self.refresh_list()
        # Keep the selection on the edited item
        self.model_listbox.selection_set(self.current_index)
        self.model_listbox.see(self.current_index)
        self.on_select(None) # Manually trigger to update status label
        
        # message box can be intrusive, just update label
        self.lbl_status.config(text=f"{status_msg} Don't forget to Save to Disk!", bootstyle="warning")

    def save_to_file(self):
        try:
            with open(CONFIG_PATH, 'w') as f:
                json.dump(self.config_data, f, indent=2)
            messagebox.showinfo("Saved", f"Configuration saved to {CONFIG_PATH}")
            self.lbl_status.config(text=f"Saved to {CONFIG_PATH}", bootstyle="success")
        except Exception as e:
            messagebox.showerror("Error", f"Failed to save file: {e}")

def main():
    # ttkbootstrap Window replaces tk.Tk()
    # themes: cosmo, flatly, journal, litera, lumen, minty, pulse, sandstone, united, yeti
    # dark themes: cyborg, darkly, solar, superhero
    root = ttk.Window(themename="cosmo") 
    
    app = ConfigEditorApp(root)
    root.mainloop()

if __name__ == "__main__":
    main()