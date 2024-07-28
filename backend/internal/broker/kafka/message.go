package kafka

const PollTimeout = 100

//func (k *Kafka) SaveConsumer(ctx context.Context, number int) error {
//	const op = scope + "SaveConsumer"
//	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
//		"bootstrap.servers":               k.cfg.Servers,
//		"group.id":                        "message_handlers",
//		"go.application.rebalance.enable": k.cfg.Rebalance,
//	})
//	defer consumer.Close()
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//	err = consumer.Subscribe(topics.HandleMessage, nil)
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//	//k.log.With()
//	//MIN_COMMIT_COUNT
//	var mtx sync.Mutex
//	go func(mtx *sync.Mutex) {
//		defer mtx.Unlock()
//		msgCount := 0
//		for {
//			select {
//			case <-ctx.Done():
//				return
//			//case sig := <-sigchan:
//			//	k.log.Error(fmt.Sprintf("Caught signal %v: terminating", sig), sl.Op(op), sl.Err(err))
//			//	return
//			default:
//				ev := consumer.Poll(PollTimeout)
//				switch e := ev.(type) {
//				case kafka.AssignedPartitions:
//					if err = consumer.Assign(e.Partitions); err != nil {
//						k.log.Error("couldn't assign partition", sl.Op(op), sl.Err(err))
//					}
//				case kafka.RevokedPartitions:
//					if err = consumer.Unassign(); err != nil {
//						k.log.Error("couldn't unassign partition", sl.Op(op), sl.Err(err))
//					}
//				case *kafka.Message:
//					msgCount += 1
//					if msgCount%k.cfg.Consumers == number {
//						if _, err = consumer.Commit(); err != nil {
//							k.log.Error("couldn't commit msg", sl.Op(op), sl.Err(err))
//						}
//					}
//					var msg api.Message
//					if err = json.Unmarshal(e.Value, &msg); err != nil {
//						k.log.Error("couldn't unmarshal msg", sl.Op(op), sl.Err(err))
//					}
//					id, err := k.storage.Save(ctx, &msg)
//					if err != nil {
//						k.log.Error("couldn't save msg", sl.Op(op), sl.Err(err))
//					}
//
//				case kafka.PartitionEOF:
//					k.log.Error("reached end of partition", sl.Op(op))
//				case kafka.Error:
//					k.log.Error("got error of consumer", sl.Op(op), sl.Err(ev.(kafka.Error)))
//					return
//				default:
//					k.log.Error("Ignored")
//				}
//			}
//		}
//	}(&mtx)
//	mtx.Lock()
//
//	return nil
//}
