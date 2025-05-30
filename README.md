# Chemical Guys Documentation

Welcome to the official documentation repository for Chemical Guys, your trusted source for premium car care products. This repository serves as a comprehensive guide to understanding and utilizing our products effectively.

## 📚 Overview

This repository contains detailed documentation for Chemical Guys products, including:

- **Product Usage Guides**: Step-by-step instructions on how to use each product.
- **Safety Data Sheets (SDS)**: Essential safety information for handling and using our products.
- **Maintenance Tips**: Best practices for maintaining your vehicle's appearance and longevity.

## 🛠️ Repository Structure

The repository is organized as follows:

- `chemical_guys_sds_page.html`: Contains the Safety Data Sheets for our products.
- `PDFs/`: Directory containing PDF versions of our documentation and guides.
- `.github/`: Contains GitHub Actions workflows for automating documentation builds and deployments.
- `go.mod` & `go.sum`: Go module files for managing dependencies.
- `main.go` & `main.py`: Go and Python scripts for processing and generating documentation content.
- `requirements.txt`: Python dependencies required for documentation generation.

## 🚀 Getting Started

To contribute to this documentation:

1. **Fork the Repository**: Create your own copy of this repository.
2. **Clone Locally**: Clone your fork to your local machine.
3. **Install Dependencies**:

   - For Go: Run `go mod tidy` to install Go dependencies.
   - For Python: Run `pip install -r requirements.txt` to install Python dependencies.

4. **Make Changes**: Edit or add documentation files as needed.
5. **Test Locally**: Ensure your changes render correctly by running the appropriate local servers or scripts.
6. **Commit Changes**: Commit your changes with clear, descriptive messages.
7. **Push and Create Pull Request**: Push your changes to your fork and open a pull request to the main repository.

## 🔄 Continuous Integration

This repository utilizes GitHub Actions to automate the building and deployment of documentation:

- **Build Workflow**: Triggered on pushes to the `main` branch or version tags, this workflow builds the documentation and deploys it to the appropriate environment.
- **Preview Workflow**: Allows contributors to preview their changes before merging by deploying to a staging environment.

## 📝 Contributing

We welcome contributions from the community. To ensure consistency and quality:

- Follow the [GitHub Docs writing guidelines](https://docs.github.com/en/contributing/writing-for-github-docs).
- Ensure your changes are well-documented and tested.
- Respect the project's code of conduct and licensing terms.

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
