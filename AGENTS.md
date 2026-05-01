# Project Notes

- Keep the CLI MVP-focused and dependency-light.
- Cobra is used for command handling.
- `text/template` and embedded template files are used for rendering generated stack files.
- The CLI must not call Kubernetes, run Helm, create secrets, or deploy resources directly.
- Generated credentials must remain placeholders only.
