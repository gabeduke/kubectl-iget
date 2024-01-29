# Kubectl Interactive Get (kubectl-iget) Plugin

## Overview
`kubectl-iget` is an interactive Kubernetes CLI plugin that extends the functionality of `kubectl get`. It allows users to interactively construct object selectors and filters for Kubernetes objects, enhancing the experience of querying Kubernetes resources.

## Features
- **Drop-In Replacement for `kubectl get`**: Use `kubectl iget` just like you would use `kubectl get`.
- **Interactive Selection and Filtering**: Dynamically build selectors and filters for Kubernetes objects through an interactive CLI interface.
- **Support for Various Kubernetes Objects**: Works with standard Kubernetes objects and CRDs.
- **Customizable Queries**: Tailor your queries with command-line arguments and interactive selections.

## Getting Started
Install the plugin and use it as a replacement for `kubectl get`. For example:

```bash
kubectl iget pods -l app=<myapp> -n <mynamespace>
```

## Getting Started
(TODO: Installation instructions, basic usage, etc.)

## Usage
Use `kubectl iget` followed by the Kubernetes object type and any additional `kubectl get` flags you normally use. The interactive interface will then guide you through further refining your query.

### Basic Example

```bash
kubectl iget pods -l app=myapp -n mynamespace
```

### Advanced Features
(TODO: Describe advanced features like the dry-run option and potential future features.)

## Contributing
Contributions to the project are welcome! Please refer to the contributing guidelines for more information.

## License
(TODO: Specify the license.)

## Acknowledgements
(TODO: Acknowledge any contributors, inspirations, or significant dependencies.)
