# Project Context: Droid Config (Factory Config Editor)

## Project Overview

**Droid Config** (also referred to as "Factory Config Editor") is a lightweight Python GUI application designed to manage configuration files for what appears to be a larger system (possibly related to AI model configurations, given the fields like `api_key`, `base_url`, `provider`).

The specific target file is `~/.factory/config.json`. The tool allows users to create, read, update, and delete model configurations stored in this JSON file.

## Tech Stack

*   **Language:** Python 3.6+
*   **GUI Framework:** `CustomTkinter` (Modern, Rounded UI)
*   **Data Format:** JSON
*   **Packaging:** `setuptools`

## Architecture & Key Files

The project follows a standard Python package structure:

*   **`src/droid_config/main.py`**: The core application logic. It initializes a `ttkbootstrap.Window` and uses bootstrap-styled widgets.
*   **`setup.py`**: Configuration for building and installing the package. Includes `ttkbootstrap` as a dependency.
*   **`pyproject.toml`**: Standard build system requirement definitions.

## Installation & Usage

### Build and Install

To install the package and its dependencies (including `ttkbootstrap`):

```bash
pip install .
```

### Running the Application

After installation, the application can be launched from the command line:

```bash
droid-config
```

*Note: The `README.md` mentions `factory-editor` as the command, but `setup.py` registers `droid-config`. Check which one is active or if `factory-editor` is intended to be an alias.*

## Development Conventions

*   **Style:** The code uses standard Python formatting.
*   **UI/UX:** The interface is built with `ttk` widgets for a somewhat native look. It uses a PanedWindow layout with a list on the left and a form on the right.
*   **Configuration:** The app defaults to `~/.factory/config.json`. Ensure this directory exists or the app has permissions to create it (the code attempts to create it).
