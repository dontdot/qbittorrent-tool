from setuptools import setup, find_packages

setup(
    name="qbittorrent-tool",
    version="1.0.0",
    author="fengqi",
    description="A tool for managing qBittorrent",
    packages=find_packages(),
    install_requires=[
        "requests>=2.25.0",
    ],
    entry_points={
        'console_scripts': [
            'qbittorrent-tool=main:main',
        ],
    },
    python_requires='>=3.6',
)