package todo

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type Todohandler struct {
	repository *TodoRepository
}

func (handler *Todohandler) GetAll(c *fiber.Ctx) error {
	var todos []Todo = handler.repository.FindAll()
	return c.JSON(todos)
}

func (handler *Todohandler) Get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	todo, err := handler.repository.Find(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status": 404,
			"error":  err,
		})
	}

	return c.JSON(todo)
}

func (handler *Todohandler) Create(c *fiber.Ctx) error {
	data := new(Todo)

	if err := c.BodyParser(data); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"error":   err,
		})
	}

	item, err := handler.repository.Create(*data)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "failed creating item",
			"error":   err,
		})
	}

	return c.JSON(item)
}

func (handler *Todohandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Item not found",
			"error":   err,
		})
	}

	todo, err := handler.repository.Find(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "item not found",
		})
	}

	todoData := new(Todo)

	if err := c.BodyParser(todoData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "review your input",
			"error":   err,
		})
	}

	todo.Name = todoData.Name
	todo.Description = todoData.Description
	todo.Status = todoData.Status

	item, err := handler.repository.Save(todo)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error updating todo",
			"error":   err,
		})
	}

	return c.JSON(item)
}

func (handler *Todohandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "failed to detele todo",
			"error":   err,
		})
	}
	RowsAffected := handler.repository.Delete(id)
	statusCode := 204
	if RowsAffected == 0 {
		statusCode = 400
	}

	return c.Status(statusCode).JSON(nil)
}

func NewTodohandler(repository *TodoRepository) *Todohandler {
	return &Todohandler{
		repository: repository,
	}
}

func Register(router fiber.Router, database *gorm.DB) {
	database.AutoMigrate(&Todo{})
	todoRepository := NewTodoRepository(database)
	todoHandler := NewTodohandler(todoRepository)

	movieRouter := router.Group("/todo")
	movieRouter.Get("/", todoHandler.GetAll)
	movieRouter.Get("/:id", todoHandler.Get)
	movieRouter.Put("/:id", todoHandler.Update)
	movieRouter.Post("/", todoHandler.Create)
	movieRouter.Delete("/:id", todoHandler.Delete)

}
