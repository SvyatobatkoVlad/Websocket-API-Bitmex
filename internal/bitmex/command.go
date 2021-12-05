package bitmex

const (
	// Instrument is a topic to subscribe
	Instrument Topic = "instrument"
)

type (
	// Command is command layout
	Commands struct {
		Action  string   `json:"action"`
		Symbols []string `json:"symbols,omitempty"`
	}
	// ErrorResponse is error response layout
	ErrorResponse struct {
		Code    string `json:"error"`
		Message string `json:"message,omitempty"`
	}
	// Topic is for topics to subscribe
	Topic string
)

func CommandExecution(w *WebsocketClient, command *Commands) error {
	bitmexCommand := Command{
		Op:   command.Action,
		Args: command.Symbols,
	}
	if len(bitmexCommand.Args) != 0 {
		for index, arg := range bitmexCommand.Args {
			bitmexCommand.Args[index] = string(Instrument) + ":" + arg
		}
	} else {
		bitmexCommand.Args = []string{
			string(Instrument),
		}
	}

	err := w.SendCommand(bitmexCommand)
	if err != nil {
		w.logger.Warning(err)
		return err
	}

	return nil
}
