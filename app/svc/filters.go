package svc

//
//func FilterDeltaFunc(back *PriceSvc, deltas map[string]int, expiryPeriod time.Duration) PriceFilterFunc {
//	return func(ctx context.Context, record *PriceRecord) bool {
//		app.Logger().Warn(ctx, "filter: %+v", record)
//		deltaPercent := deltas[record.Symbol]
//		if deltaPercent == 0 || deltaPercent >= 100 {
//			return true
//		}
//		current, err := back.GetLatestSymbolPrice(ctx, record.Symbol, record.Currency)
//		if err != nil {
//			if !errors.Is(err, ErrNoValue) {
//				app.Logger().Error(ctx, "error getting current for filter: %s", err)
//				return false
//			}
//			return true
//		}
//		if current.TimeStamp.Add(expiryPeriod).Before(time.Now()) {
//			return true
//		}
//		deltaDiff := current.Price * float64(deltaPercent) / 100
//		diff := record.Price - current.Price
//		if diff < 0 {
//			diff *= -1
//		}
//		if diff >= deltaDiff {
//			return true
//		}
//		return false
//	}
//}
