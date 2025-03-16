package animations

type Animation struct {
	First        int
	Last         int
	Step         int     // Step between frames
	SpeedInTps   float32 // Speed in ticks per second
	frameCounter float32
	frame        int
}

func (a *Animation) Update() {
	a.frameCounter -= 1.0
	if a.frameCounter < 0.0 {
		a.frameCounter = a.SpeedInTps
		a.frame += a.Step
		if a.frame > a.Last {
			a.frame = a.First
		}
	}
}

func (a *Animation) Frame() int {
	return a.frame
}

func NewAnimation(first, last, step int, speedInTps float32) *Animation {
	return &Animation{
		First:        first,
		Last:         last,
		Step:         step,
		SpeedInTps:   speedInTps,
		frameCounter: speedInTps,
		frame:        first,
	}
}
