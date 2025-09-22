env "local" {
  src = "file://schema"
  dev = "docker://mysql/8/oshiapp"
  url = "mysql://oshiapp_user:oshiapp_password@localhost:3306/oshiapp"
}
