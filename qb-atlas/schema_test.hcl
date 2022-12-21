schema "public" {
	charset = "UTF-8"
}

table "users" {
	schema = schema.public
	column "id" {
		type = serial
	}

	column "username" {
		type = varchar(100)
	}

	column "password" {
		type = varchar
	}

	column "bio" {
		type = text
		null = true
	}

	index "uq_users_username" {
		columns = [
			column.username
		]
		unique = true
	}

	primary_key {
		columns = [column.id]
	}
}

table "blogs" {
	schema = schema.public
	column "id" {
		type = serial
	}

	column "user_id" {
		type = int
	}

	column "title" {
		type = varchar
	}

	column "draft" {
		type = boolean
		default = false
	}

	column "created_at" {
		type = timestamp
		default = "NOW()"
	}

	column "content" {
		type = text
	}

	foreign_key "fk_blogs_user_id" {
		columns = [column.user_id]
		ref_columns = [table.users.column.id]
	}

	primary_key {
		columns = [column.id]
	}
}
