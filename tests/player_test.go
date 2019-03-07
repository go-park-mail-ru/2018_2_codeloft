package tests

import (
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2018_2_codeloft/game/game/models"
)

func TestChangeDirectionUnknown(t *testing.T) {
	inPlayer := models.Player{MoveDirection: "LEFT"}
	expectedPlayer := models.Player{MoveDirection: "LEFT"}

	inPlayer.ChangeDirection("WRONG")
	if !reflect.DeepEqual(inPlayer, expectedPlayer) {
		t.Errorf("Unknown direction should change nothing Expected to be %v but got %v", inPlayer, expectedPlayer)
	}
}

func TestChangeDirectionLeftToRight(t *testing.T) {
	inPlayer := models.Player{MoveDirection: "LEFT"}
	expectedPlayer := models.Player{MoveDirection: "LEFT"}

	inPlayer.ChangeDirection("RIGHT")
	if !reflect.DeepEqual(inPlayer, expectedPlayer) {
		t.Errorf("Opposite direction should change nothing Expected to be %v but got %v", inPlayer, expectedPlayer)
	}
}

func TestChangeDirectionRightToLeft(t *testing.T) {
	inPlayer := models.Player{MoveDirection: "RIGHT"}
	expectedPlayer := models.Player{MoveDirection: "RIGHT"}

	inPlayer.ChangeDirection("LEFT")
	if !reflect.DeepEqual(inPlayer, expectedPlayer) {
		t.Errorf("Opposite direction should change nothing  Expected to be %v but got %v", inPlayer, expectedPlayer)
	}
}

func TestChangeDirectionUpToDown(t *testing.T) {
	inPlayer := models.Player{MoveDirection: "UP"}
	expectedPlayer := models.Player{MoveDirection: "UP"}

	inPlayer.ChangeDirection("DOWN")
	if !reflect.DeepEqual(inPlayer, expectedPlayer) {
		t.Errorf("Opposite direction should change nothing  Expected to be %v but got %v", inPlayer, expectedPlayer)
	}
}

func TestChangeDirectionDownToUp(t *testing.T) {
	inPlayer := models.Player{MoveDirection: "DOWN"}
	expectedPlayer := models.Player{MoveDirection: "DOWN"}

	inPlayer.ChangeDirection("UP")
	if !reflect.DeepEqual(inPlayer, expectedPlayer) {
		t.Errorf("Opposite direction should change nothing Expected to be %v but got %v", inPlayer, expectedPlayer)
	}
}

func TestChangeDirectionOk(t *testing.T) {
	inPlayer := models.Player{MoveDirection: "LEFT"}
	expectedPlayer := models.Player{MoveDirection: "UP"}

	inPlayer.ChangeDirection("UP")
	if !reflect.DeepEqual(inPlayer, expectedPlayer) {
		t.Errorf("Direction should be changed as attribute Expected to be %v but got %v", inPlayer, expectedPlayer)
	}
}

func TestMoveLeftBorder(t *testing.T) {
	inPlayer := models.Player{MoveDirection: "LEFT", Position: models.Position{0, 0}}
	expectedPlayer := models.Player{MoveDirection: "LEFT", Position: models.Position{models.FIELD_WIDTH - 1, 0}, Tracer: []models.Position{models.Position{models.FIELD_WIDTH - 1, 0}}}

	inPlayer.Move()
	if !reflect.DeepEqual(inPlayer, expectedPlayer) {
		t.Errorf("Position should be set to right border Expected to be %v but got %v", inPlayer, expectedPlayer)
	}
}

func TestMoveRightBorder(t *testing.T) {
	inPlayer := models.Player{MoveDirection: "RIGHT", Position: models.Position{models.FIELD_WIDTH - 2, 0}}
	expectedPlayer := models.Player{MoveDirection: "RIGHT", Position: models.Position{0, 0}, Tracer: []models.Position{models.Position{0, 0}}}

	inPlayer.Move()
	if !reflect.DeepEqual(inPlayer, expectedPlayer) {
		t.Errorf("Position should be set to left border Expected to be %v but got %v", inPlayer, expectedPlayer)
	}
}

func TestMoveUpBorder(t *testing.T) {
	inPlayer := models.Player{MoveDirection: "UP", Position: models.Position{0, 0}}
	expectedPlayer := models.Player{MoveDirection: "UP", Position: models.Position{0, models.FIELD_HEIGHT - 1}, Tracer: []models.Position{models.Position{0, models.FIELD_HEIGHT - 1}}}

	inPlayer.Move()
	if !reflect.DeepEqual(inPlayer, expectedPlayer) {
		t.Errorf("Position should be set to down border Expected to be %v but got %v", inPlayer, expectedPlayer)
	}
}

func TestMoveDownBorder(t *testing.T) {
	inPlayer := models.Player{MoveDirection: "DOWN", Position: models.Position{0, models.FIELD_HEIGHT - 2}}
	expectedPlayer := models.Player{MoveDirection: "DOWN", Position: models.Position{0, 0}, Tracer: []models.Position{models.Position{0, 0}}}

	inPlayer.Move()
	if !reflect.DeepEqual(inPlayer, expectedPlayer) {
		t.Errorf("Position should be set to Up border Expected to be %v but got %v", inPlayer, expectedPlayer)
	}
}

func TestMoveOk(t *testing.T) {
	originTrace := []models.Position{models.Position{50, 30}}
	inPlayer := models.Player{MoveDirection: "DOWN", Position: models.Position{50, 30}, Tracer: originTrace}
	expectedPlayer := models.Player{MoveDirection: "DOWN", Position: models.Position{50, 31}, Tracer: append(originTrace, models.Position{50, 31})}

	inPlayer.Move()
	if !reflect.DeepEqual(inPlayer, expectedPlayer) {
		t.Errorf("Expected to be %v but got %v", inPlayer, expectedPlayer)
	}
}
