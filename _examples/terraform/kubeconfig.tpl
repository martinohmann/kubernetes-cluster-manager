apiVersion: v1
clusters:
- cluster:
    server: ${endpoint}
    certificate-authority-data: ${certificate_authority_data}
  name: ${clustername}
contexts:
- context:
    cluster: ${clustername}
    user: aws
  name: ${context}
current-context: ${context}
kind: Config
preferences: {}
users:
- name: aws
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      command: aws-iam-authenticator
      args:
        - "token"
        - "-i"
        - "${clustername}"
