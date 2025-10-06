# Steps to using GORM
## Impact the database

1. Open `internal/adapters/storage/database/models`
2. Create a new file with the name of the table to be created. Example `new_table.go`
3. Define the new model fields

```go
package gorm_models
   
import "gorm.io/gorm"

// NewTable represents a new table of the database.
type NewTable struct {
    gorm.Model
    Name string `gorm:"size:255"`
}
```

4. Open `internal/adapters/storage/database/db.go` and add the new model to the existing models list. This function will be impact at Database.

```go
package postgres

// {...} - Rest of logic

// Executes migrations and impacts the database
func ExecMigrations() {
	migrate(
		&gorm_models.User{},
		&gorm_models.NewGormModel{},
	)
}
```

5. Open ``cmd/http`` and run the next command: `go build && ./http.exe` or run debbugger in the `main.go` file
6. Verify the new changes in the database using `pgAdmin` or other database manager (GUI)

## Loading tables relationship

When we creates a new model in `internal/database/gorm_models`, and this model has a relation with other table. For example `User <--> []AuthAccount`


```go
package models
   
import "gorm.io/gorm"

type User struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
    FullName  string    `gorm:"size:255;default:null"`
    ...

    // Foreign key --> AuthAccount Table
    AuthAccounts []AuthAccount `gorm:"foreignKey:UserID"`
}
```
```go
package models
   
import "gorm.io/gorm"

type AuthAccount struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Provider      string    `gorm:"not null"`
    ...

    // Foreign key --> User Table
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
```

And we call a repository to obtain User data saved in database:
```go
package repository

func (repo *RepoConnection) FindUserByID(id string) (*gorm_models.User, error) {
	var user gorm_models.User

	if result := repo.db.First(&user, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
```


GORM doesn't automatically load the relations of the both tables, so in this case `AuthAccount` data will be empty `nil` or `len = 0` and the code throught an `error`.

To solve this problem we have to use `Preload()` function before search the first result in database. This function will automatically load the relations and we will have all the necessary data available.

````go
func (repo *RepoConnection) FindUserByID(id string) (*gorm_models.User, error) {
	var user gorm_models.User

	if result := repo.db.Preload(string(AuthAccountKey)).First(&user, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
````
