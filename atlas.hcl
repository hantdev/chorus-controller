env "local" {
  src = "file://migrations"
  dev = "docker://postgres/16/dev?search_path=public"
  url = "postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable"
}

env "dev" {
  src = "file://migrations"
  dev = "docker://postgres/16/dev?search_path=public"
  url = "postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable"
}
