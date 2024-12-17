package repollo

import (
	"errors"
	"fmt"
	"sync"
)

// Event representa un cambio en la colección.
type Event[T any] struct {
	Event string
	Key   string
	Value T
}

// Collection es una estructura genérica que representa una colección de documentos.
type Collection[T any] struct {
	mu     sync.Mutex
	data   map[string]T
	events chan Event[T]
	closed bool
	wg     sync.WaitGroup // WaitGroup para rastrear las operaciones concurrentes.
}

// NewCollection crea una nueva colección con un canal de eventos.
func NewCollection[T any]() *Collection[T] {
	return &Collection[T]{
		data:   make(map[string]T),
		events: make(chan Event[T], 100), // Canal con buffer para manejar eventos.
	}
}

// Create crea un nuevo documento en la colección.
func (c *Collection[T]) Create(key string, value T) (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return value, errors.New("collection is closed")
	}

	if _, exists := c.data[key]; exists {
		return value, errors.New("document already exists")
	}

	c.data[key] = value
	c.publishEvent("create", key, value)
	return value, nil
}

// Get lee un documento de la colección.
func (c *Collection[T]) Get(key string) (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero T
	if value, exists := c.data[key]; exists {
		return value, nil
	}

	return zero, errors.New("document not found")
}

// Update actualiza un documento existente en la colección.
func (c *Collection[T]) Update(key string, value T) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.New("collection is closed")
	}

	if _, exists := c.data[key]; !exists {
		return errors.New("document not found")
	}

	c.data[key] = value
	c.publishEvent("update", key, value)
	return nil
}

// Delete elimina un documento de la colección.
func (c *Collection[T]) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.New("collection is closed")
	}

	value, exists := c.data[key]
	if !exists {
		return errors.New("document not found")
	}

	delete(c.data, key)
	c.publishEvent("delete", key, value)
	return nil
}

func (c *Collection[T]) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.data)
}

// Where aplica un filtro basado en un predicado.
func (c *Collection[T]) Where(predicate func(T) bool) *QueryResult[T] {
	c.mu.Lock()
	defer c.mu.Unlock()

	var filtered []T
	for _, value := range c.data {
		if predicate(value) {
			filtered = append(filtered, value)
		}
	}

	return &QueryResult[T]{results: filtered}
}

func (c *Collection[T]) Find(predicate func(T) bool) (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero T
	for _, value := range c.data {
		if predicate(value) {
			return value, nil
		}
	}

	return zero, errors.New("document not found")
}

// Events devuelve el canal de eventos para que otros puedan leerlo.
func (c *Collection[T]) Events() <-chan Event[T] {
	return c.events
}

// Close cierra el canal de eventos. Espera a que todas las operaciones concurrentes terminen antes de cerrar.
func (c *Collection[T]) Close() {
	c.wg.Wait() // Esperar a que todas las operaciones concurrentes terminen.
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		c.closed = true
		close(c.events)
	}
}

// publishEvent publica un evento en el canal.
func (c *Collection[T]) publishEvent(event string, key string, value T) {
	c.wg.Add(1) // Incrementar el contador del WaitGroup.
	go func() {
		defer c.wg.Done() // Decrementar el contador al terminar.
		select {
		case c.events <- Event[T]{Event: event, Key: key, Value: value}:
			// Evento enviado correctamente.
		default:
			fmt.Println("Warning: Event channel is full, event dropped")
		}
	}()
}
