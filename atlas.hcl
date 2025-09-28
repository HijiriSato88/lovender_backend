env "local" {
  src = "file://schema"
  dev = "docker://mysql/8/lovender"
  url = "mysql://lovender_user:lovender_password@localhost:3306/lovender"
}

env "production" {
  src = "file://schema"
  dev = "docker://mysql/8/lovender"
  url = "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST})/${DB_NAME}"
}