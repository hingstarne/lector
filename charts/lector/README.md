# lector helm chart

## Prerequisites

1. First ensure that the secrets can be decrypted from the key in the *main* AWS account.
```bash
export AWS_PROFILE=fl-current
export AWS_REGION=eu-west-1
```

1. Make sure you have your `~/.aws/credentials` exist. But following the instructions here:
https://github.com/hingstarne/fl-cli-tools/blob/master/README.md#prerequisites

1. Install **helm**:
https://github.com/kubernetes/helm#install

*On MacOS*
```
brew install kubernetes-helm
```

1. Install [Helm secrets plugin](https://github.com/futuresimple/helm-secrets)
```
helm plugin install https://github.com/futuresimple/helm-secrets
```

## Updating secrets

You can then edit secrets with SOPS via:
```bash
helm secrets edit secrets.yaml
```

## Permission errors
If you don't have permission to access this key you will see an error;

```bash
Failed to get the data key required to decrypt the SOPS file.
```

