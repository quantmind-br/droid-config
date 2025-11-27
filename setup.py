from setuptools import setup, find_packages

setup(
    name="droid-config",
    version="0.1.0",
    description="GUI editor for ~/.factory/config.json",
    author="Gemini CLI",
    packages=find_packages(where="src"),
    package_dir={"": "src"},
    entry_points={
        "console_scripts": [
            "droid-config=droid_config.main:main",
        ],
    },
    python_requires=">=3.6",
    install_requires=[
        "ttkbootstrap",
    ],
)