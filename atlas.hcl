env "local" {
  src = [
    "file://DB/schema/schema.sql",
    "file://DB/migrations/002_foreign_keys.up.sql"
  ]
  dev = "postgres://admin:password@db:5432/dev?sslmode=disable&search_path=public"
}
