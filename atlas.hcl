data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlas"
  ]
}

env "local" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev?search_path=public"
  url = "postgres://fazt:mysecretpassword@localhost:5432/tiny_store?search_path=public&sslmode=disable"
  migration {
    dir = "file://migrations"
    revisions_schema = "public"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }

}

env "production" {
  src = data.external_schema.gorm.url
  url = getenv(DB_DATABASE_URL)
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}