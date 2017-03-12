package tic_tac_toe_minimax

import "math"

type StateTreeNode struct {
	Board      BoardState
	Children   []StateTreeNode
	NextPlayer Player
	Fitness    int32
	BestChild  int
	OurPlayer  Player
}

func NewDefaultState(us Player) StateTreeNode {
	return StateTreeNode{
		Board:      NewBoard(3),
		Children:   []StateTreeNode{},
		NextPlayer: PlayerX,
		Fitness:    0,
		OurPlayer:  us,
	}
}

func NewState(size int, us Player) StateTreeNode {
	return StateTreeNode{
		Board:      NewBoard(size),
		Children:   []StateTreeNode{},
		NextPlayer: PlayerX,
		Fitness:    0,
		OurPlayer:  us,
	}
}

func (s *StateTreeNode) FindAllChildStates() {
	if (s.Board.Full()) {
		s.Fitness = 0
		return
	}
	if p := s.Board.CheckWin(); p != NoPlayer {
		if p == s.OurPlayer {
			// we win, therefore state is highly desirable
			s.Fitness = math.MaxInt32
		} else {
			// we don't win, therefore state is not desirable
			s.Fitness = -1
		}
		return
	}

	// calculate child states
	moves := s.Board.OpenCells()
	s.Children = make([]StateTreeNode, len(moves))
	for i, pos := range moves {
		s.Children[i].Board = CopyBoard(&s.Board)
		s.Children[i].Board.SetPos(pos, s.NextPlayer)
		s.Children[i].NextPlayer = OppositePlayer(s.NextPlayer)
		s.Children[i].OurPlayer = s.OurPlayer
		s.Children[i].FindAllChildStates()
	}

	// decide what our fitness is, based on children
	if s.NextPlayer == s.OurPlayer {
		// we are playing next, so do max of children
		s.Fitness = math.MinInt32
		for i, child := range s.Children {
			if s.Fitness < child.Fitness {
				s.Fitness = child.Fitness
				s.BestChild = i
			}
		}
	} else {
		// other player playing, minimize children
		s.Fitness = math.MaxInt32
		for i, child := range s.Children {
			if s.Fitness > child.Fitness {
				s.Fitness = child.Fitness
				s.BestChild = i
			}
		}
	}
}

func (s *StateTreeNode) GetChildForMove(pos Pos) StateTreeNode {
	for _, child := range s.Children {
		if child.Board.GetPos(pos) != NoPlayer {
			return child
		}
	}
	// todo
	return s.Children[0]
}

func OppositePlayer(player Player) Player {
	switch player {
	case PlayerX:
		return PlayerO
	case PlayerO:
		return PlayerX
	default:
		return PlayerX
	}
}

func max(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}
func min(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}
