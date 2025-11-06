# Deepgram Self-Hosted Helm Chart (Voiceflow Fork)

This is Voiceflow's fork of the official [Deepgram Self-Hosted Resources](https://github.com/deepgram/self-hosted-resources).

## Why We Forked

We maintain our own version of the Deepgram Helm chart to support Voiceflow-specific requirements:

- **Datadog Integration**: Added support for Datadog service naming via `DD_SERVICE` environment variables for better observability
- **Custom Environment Variables**: Enhanced templating to allow flexible environment variable configuration through `values.yaml`
- **Internal Deployment Needs**: Tailored configurations for Voiceflow's infrastructure and deployment patterns

## Versioning

We use a hybrid versioning scheme: `<upstream-version>-vf.<voiceflow-version>`

Example: `0.23.1-vf.1` means:
- Based on upstream version `0.23.1`
- Voiceflow-specific changes version `1`

## Contents

* [Helm Chart](charts/deepgram-self-hosted/README.md) for Kubernetes deployments
* [Docker Compose Files](./docker/README.md) for deploying with Docker
* [Podman Compose Files](./podman/README.md) for deploying with Podman
* [Diagnostic](./diagnostics/README.md) tools and scripts for troubleshooting deployments

## Documentation

You can learn more about the Deepgram API at [developers.deepgram.com](https://developers.deepgram.com/docs), and more about self-hosting Deepgram in the [relevant documentation](https://developers.deepgram.com/docs/self-hosted-introduction).

## Getting Help

We love to hear from you so if you have questions, comments or find a bug in this repo, let us know!

- If you have a Premium or VIP Support Plan with Deepgram you can find details and links to contact us for support on your [Console dashboard](https://console.deepgram.com).
- If you're interested in learning more about self-hosting Deepgram products, [contact us here](https://deepgram.com/contact-us)!
- If you have a specific bug or feature request for these resources, [you can open an issue in this repo](https://github.com/deepgram/self-hosted-resources/issues/new/choose).
- The [Deepgram documentation](https://developers.deepgram.com) and [Deepgram Help Center](https://deepgram.gitbook.io/help-center) have answers to many common questions.
- Deepgram also has a [developer community](https://community.deepgram.com/)!

