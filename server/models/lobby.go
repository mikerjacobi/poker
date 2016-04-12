package models

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func JoinGame(db *mgo.Database, gameID string, accountID string, buyinAmount int) (Game, error) {
	games := db.C("games")
	g := Game{}
	query := bson.M{"gameID": gameID}
	if err := games.Find(query).One(&g); err != nil {
		return Game{}, err
	}
	for i := range g.Players {
		if g.Players[i].AccountID == accountID {
			return g, nil
		}
	}

	//subtract buyin amount from account
	if buyinAmount <= 0 {
		return g, fmt.Errorf("buyin must be >0")
	}
	account, err := LoadAccountByID(db, accountID)
	if err != nil {
		return g, fmt.Errorf("failed to load account")
	}
	if account.Balance-buyinAmount < 0 {
		return g, fmt.Errorf("not enough funds to buy in")
	}
	account.Balance -= buyinAmount
	if err := account.Update(db); err != nil {
		return g, fmt.Errorf("failed to update account balance")
	}

	gp := GamePlayer{AccountID: account.AccountID, Name: account.Username, Chips: buyinAmount}
	g.Players = append(g.Players, gp)
	if err := games.Update(query, g); err != nil {
		return Game{}, err
	}
	return g, nil
}

func LeaveGame(db *mgo.Database, gameID string, accountID string) (Game, error) {
	games := db.C("games")
	g := Game{}
	query := bson.M{"gameID": gameID}
	if err := games.Find(query).One(&g); err != nil {
		return Game{}, err
	}
	for i := range g.Players {
		if g.Players[i].AccountID != accountID {
			continue
		}

		if g.GameType == "highcard" {
			if err := highCardFoldOutOfTurn(gameID, accountID); err != nil {
				return Game{}, fmt.Errorf("failed to fold highcard out of turn: %+v", err)
			}
		} else {
			return Game{}, fmt.Errorf("game %s fold out of turn not implemented.")
		}

		chips := g.Players[i].Chips
		g.Players = append(g.Players[0:i], g.Players[i+1:]...)
		if err := games.Update(query, g); err != nil {
			return Game{}, err
		}

		account, err := LoadAccountByID(db, accountID)
		if err != nil {
			return g, fmt.Errorf("failed to load account")
		}
		account.Balance += chips
		if err := account.Update(db); err != nil {
			return g, fmt.Errorf("failed to update account balance")
		}

		break

	}
	return g, nil
}

func RemovePlayerFromGames(db *mgo.Database, accountID string) error {
	gamesDB := db.C("games")
	games := []Game{}
	query := bson.M{}
	if err := gamesDB.Find(query).All(&games); err != nil {
		return err
	}

	for _, game := range games {
		for i := range game.Players {
			if game.Players[i].AccountID == accountID {
				game.Players = append(game.Players[0:i], game.Players[i+1:]...)
				if err := gamesDB.Update(bson.M{"gameID": game.ID}, game); err != nil {
					return err
				}
				break
			}
		}

	}
	return nil
}
