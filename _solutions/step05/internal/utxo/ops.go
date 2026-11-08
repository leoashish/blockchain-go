package utxo

import (
	"sort"
)

// REFERENCE SOLUTION for Step 05. Try it yourself first!

func (l *Ledger) Coinbase(id, owner string, amount int64) error {
	if amount <= 0 {
		return ErrBadAmount
	}
	if id == "" {
		return ErrBadAmount
	}
	if _, exists := l.unspent[id]; exists {
		return ErrDuplicateOutput
	}
	l.unspent[id] = Output{ID: id, Owner: owner, Amount: amount}
	return nil
}

func (l *Ledger) Outputs(owner string) []Output {
	out := []Output{}
	for _, o := range l.unspent {
		if o.Owner == owner {
			out = append(out, o)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (l *Ledger) Balance(owner string) int64 {
	var total int64
	for _, o := range l.unspent {
		if o.Owner == owner {
			total += o.Amount
		}
	}
	return total
}

func (l *Ledger) Transfer(id, from, to string, amount int64) error {
	if amount <= 0 {
		return ErrBadAmount
	}
	if _, exists := l.unspent[id+":0"]; exists {
		return ErrDuplicateOutput
	}
	if _, exists := l.unspent[id+":1"]; exists {
		return ErrDuplicateOutput
	}

	// Select `from`'s outputs ascending by (amount, id) until they cover the amount.
	candidates := []Output{}
	for _, o := range l.unspent {
		if o.Owner == from {
			candidates = append(candidates, o)
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Amount == candidates[j].Amount {
			return candidates[i].ID < candidates[j].ID
		}
		return candidates[i].Amount < candidates[j].Amount
	})

	var covered int64
	var selected []Output
	for _, o := range candidates {
		selected = append(selected, o)
		covered += o.Amount
		if covered >= amount {
			break
		}
	}
	if covered < amount {
		return ErrInsufficientFunds // nothing consumed yet
	}

	for _, o := range selected {
		delete(l.unspent, o.ID)
	}
	l.unspent[id+":0"] = Output{ID: id + ":0", Owner: to, Amount: amount}
	if change := covered - amount; change > 0 {
		l.unspent[id+":1"] = Output{ID: id + ":1", Owner: from, Amount: change}
	}
	return nil
}

func (l *Ledger) Spend(id string) (Output, error) {
	o, exists := l.unspent[id]
	if !exists {
		return Output{}, ErrOutputNotFound
	}
	delete(l.unspent, id)
	return o, nil
}
