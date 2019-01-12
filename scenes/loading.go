package scenes

//
//type loadingScene struct {
//	ticker         *time.Ticker
//	dots           string
//	doneLoading    chan struct{}
//	loadingMessage *text.Text
//}
//
//func (s *loadingScene) Render(canvas *pixelgl.Canvas, dt float64) {
//	select {
//	case <-s.doneLoading:
//		next(Start())
//		return
//	default:
//		select {
//		case <-s.ticker.C:
//			s.dots += "."
//			if len(s.dots) > 3 {
//				s.dots = ""
//			}
//		default:
//		}
//		canvas.Clear(colornames.Black)
//		s.loadingMessage.Clear()
//		_, _ = s.loadingMessage.WriteString("Loading" + s.dots)
//		matrix := pixel.IM.Moved(pixel.ZV.Sub(s.loadingMessage.Bounds().Center()))
//		s.loadingMessage.Draw(canvas, matrix)
//	}
//}
//
//func (s *loadingScene) Destroy() {
//	s.ticker.Stop()
//}
//
//func LoadingScene() Scene {
//	doneLoading := make(chan struct{})
//	go func() {
//		assets.Init()
//		time.Sleep(5 * time.Second)
//		close(doneLoading)
//	}()
//	var atlas = text.NewAtlas(
//		basicfont.Face7x13,
//		text.ASCII,
//	)
//	loadingMessage := text.New(pixel.V(0, 0), atlas)
//	dotTicker := time.NewTicker(time.Second)
//	return &loadingScene{
//		ticker:         dotTicker,
//		doneLoading:    doneLoading,
//		loadingMessage: loadingMessage,
//	}
//}
