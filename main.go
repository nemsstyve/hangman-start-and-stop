package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type DernierePartie struct {
	Mot            string
	Letters_found  string
	letters_tried  string
	Position_pendu int
}

func Selection_Mot(file string) string { //- 1.1 : Mise en place de l'aléatoire
	r1 := rand.New(rand.NewSource(time.Now().UnixNano())) //- 1.2 : faire sortir du fichier words
	f, _ := os.Open(file)
	b1 := make([]byte, 9999)
	n1, _ := f.Read(b1)
	var words [][]byte
	index := 0
	for indice, lettre := range string(b1[:n1]) {
		if lettre == 10 {
			words = append(words, b1[index:indice])
			index = indice + 1
		}
	}
	mot_au_hasard := string(words[r1.Intn(len(words))]) //- 1.3 : Selection aléatoire
	defer f.Close()
	return mot_au_hasard
}

func Affiche_Mot_Pendu(str string, valid []byte) {
	for i, char := range str {
		if i != 0 {
			fmt.Print(" ")
		}
		ecrit := 0
		for _, in_lv := range valid {
			if byte(char) == in_lv && ecrit == 0 {
				ecrit = 1
				fmt.Print(string(char - 'a' + 'A'))
			}
		}
		if ecrit == 0 {
			fmt.Print("_")
		}
	}
}

func Is_This_In_The_Word(letter string, mot string) bool {
	for _, char := range mot {
		if string(char) == letter {
			return true
		}
	}
	return false
}

func All_The_letters_Are_In_The_Word(letters []byte, mot string) bool {
	for _, char1 := range mot {
		ecrit := 0
		for _, char2 := range letters {
			if byte(char1) == byte(char2) {
				ecrit = 1
			}
		}
		if ecrit == 0 {
			return false
		}
	}
	return true
}

func Mot_Exact(letters []byte, mot string) bool {
	for i := range mot {
		if letters[i] != []byte(mot)[i] {
			return false
		}
	}
	return true
}

func Ending(message string, position_pendu int, mot_au_hasard string, valid_letters []byte) {
	hangman, _ := os.Open("positions/" + strconv.Itoa(position_pendu) + ".txt")
	hangman_tab := make([]byte, 9999)
	hangman_read, _ := hangman.Read(hangman_tab)
	fmt.Println(string(hangman_tab[:hangman_read]))
	Affiche_Mot_Pendu(mot_au_hasard, valid_letters)
	fmt.Println("\n" + message)
}

func main() {
	press_on := 1
	for press_on == 1 {
		mot_au_hasard := Selection_Mot("words.txt") //on lance la selection du mot aléatoire à partir de notre words.txt
		started := 0
		position_pendu := 0
		reussi := 0
		valid_letters := []byte{mot_au_hasard[len(mot_au_hasard)/2-1]}
		var letters_tried []byte
		stop := 0
		data, _ := ioutil.ReadFile("save.txt") // on implémmente le fichier de sauvegarde
		fmt.Print(len(string(data)))
		if len(os.Args) > 1 {
			if (os.Args[1]+" "+os.Args[2]) == "--startWith save.txt" && len(string(data)) > 2 { // on implémente la commande qui permettra de run l'ancien sauvegarde
				var données DernierePartie
				json.Unmarshal(data, &données)
				mot_au_hasard = données.Mot
				valid_letters = []byte(données.Letters_found)
				letters_tried = []byte(données.letters_tried)
				position_pendu = données.Position_pendu
				os.Args = os.Args[:1]
			}
		}
		for position_pendu < 10 {
			fmt.Print("\n[Attempt]\n")
			if started == 0 {
				fmt.Print("Good Luck, you have 10 attempts. (Type \"STOP\" to stop & save the game)" + "\n")
				started = 1
			}
			if position_pendu > 0 && reussi == 0 {
				fmt.Print("Not present in the word, " + strconv.Itoa(10-position_pendu) + " attempts remaining\n")
			}
			if !All_The_letters_Are_In_The_Word(valid_letters, mot_au_hasard) {
				hangman, _ := os.Open("positions/" + strconv.Itoa(position_pendu) + ".txt")
				hangman_tab := make([]byte, 9999)
				hangman_read, _ := hangman.Read(hangman_tab)
				fmt.Print(string(hangman_tab[:hangman_read]) + "\n")
				Affiche_Mot_Pendu(mot_au_hasard, valid_letters)
				fmt.Print("\n")
				fmt.Print("Already tried :")
				for _, ltr := range letters_tried {
					fmt.Print(" " + string(ltr))
				}
				fmt.Print("\nProposition : ")
				var proposition string
				fmt.Scanln(&proposition)
				if proposition == "STOP" { // ici on lui dit si le joueur tape STOP, alors le jeu se sauvegarde vers save.txt
					stop = 1
					a_sauvegarder := DernierePartie{mot_au_hasard, string(valid_letters), string(letters_tried), position_pendu}
					données, _ := json.Marshal(a_sauvegarder)
					ioutil.WriteFile("save.txt", données, 0777)
					fmt.Print("Game saved in save.txt.")
					break
				}
				if len(proposition) == 1 {
					if Is_This_In_The_Word(proposition, mot_au_hasard) || Is_This_In_The_Word(string(proposition[0]-'A'+'a'), mot_au_hasard) {
						if Is_This_In_The_Word(proposition, mot_au_hasard) {
							valid_letters = append(valid_letters, proposition[0])
						}
						if Is_This_In_The_Word(string(proposition[0]-'A'+'a'), mot_au_hasard) {
							valid_letters = append(valid_letters, proposition[0]-'A'+'a')
						}
						reussi = 1
					} else if !(Is_This_In_The_Word(proposition, string(letters_tried))) && !(Is_This_In_The_Word(string(proposition[0]-'A'+'a'), string(letters_tried))) {
						letters_tried = append(letters_tried, proposition[0])
						position_pendu++
						reussi = 0
					}
				}
			} else {
				fmt.Print("\n[End]\n")
				Ending("Congrats !", position_pendu, mot_au_hasard, valid_letters)
				break
			}
			fmt.Print("\n")
		}
		if position_pendu == 10 { //affiche les positions hangman si le mot est erroné
			fmt.Print("\n[End]\n")
			hangman, _ := os.Open("positions/10.txt")
			hangman_tab := make([]byte, 9999)
			hangman_read, _ := hangman.Read(hangman_tab)
			fmt.Print(string(hangman_tab[:hangman_read]) + "\n")
			Ending("You lose !", position_pendu, mot_au_hasard, []byte(mot_au_hasard))
		}
		if stop == 0 {
			ioutil.WriteFile("save.txt", []byte(""), 0777) //écrit dans le fichier save.txt les sauvegardes
			var start_again string
			for start_again != "n" && start_again != "N" && start_again != "o" && start_again != "O" {
				fmt.Print("\nstart_again ? (o/n) : ")
				fmt.Scanln(&start_again)
				if start_again == "n" || start_again == "N" {
					press_on = 0
				}
				if start_again == "o" || start_again == "O" {
					press_on = 1
				}
			}
		} else if stop == 1 {
			break
		}
	}
}
