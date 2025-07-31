package handlers

type History struct {
	history      []string // History of executed queries
	historyindex int      // Current index in the history
}

func (h *History) GetPreviousQuery() string {
	if h.historyindex >= 0 {
		query := h.history[h.historyindex]
		h.historyindex = max(0, h.historyindex-1)
		return query
	}
	return ""
}

func (h *History) GetNextQuery() string {
	if h.historyindex <= len(h.history)-1 {
		h.historyindex++
		if(h.historyindex>= len(h.history)) {
			h.historyindex = len(h.history) - 1 // Reset to last query if out of bounds
			return "" // Prevent going out of bounds
		}
		query := h.history[h.historyindex]
		return query
	}
	return ""
}
