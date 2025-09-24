env "local" {
  src = "file://schema"
  dev = "docker://mysql/8/lovender"
  url = "mysql://lovender_user:lovender_password@localhost:3306/lovender"
}