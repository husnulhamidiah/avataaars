package handler

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type AvatarData struct {
	AvatarStyle     []string `json:"avatarStyle"`
	AccessoriesType []string `json:"accessoriesType"`
	HairColor       []string `json:"hairColor"`
	ClotheColor     []string `json:"clotheColor"`
	FacialHairType  []string `json:"facialHairType"`
	ClotheType      []string `json:"clotheType"`
	EyeType         []string `json:"eyeType"`
	EyebrowType     []string `json:"eyebrowType"`
	MouthType       []string `json:"mouthType"`
	SkinColor       []string `json:"skinColor"`
	TopType         []string `json:"topType"`
}

const jsonData = `{
  "avatarStyle": ["Circle"],
  "accessoriesType": [
    "Blank",
    "Kurt",
    "Prescription01",
    "Prescription02",
    "Round",
    "Sunglasses",
    "Wayfarers"
  ],
  "hairColor": [
    "Auburn",
    "Black",
    "Blonde",
    "BlondeGolden",
    "Brown",
    "BrownDark",
    "PastelPink",
    "Platinum",
    "Red",
    "SilverGray"
  ],
  "clotheColor": [
    "Black", "Blue01", "Blue02", "Blue03", "Gray01", "Gray02",
    "Heather", "PastelBlue", "PastelGreen", "PastelOrange", 
    "PastelRed", "PastelYellow", "Pink", "Red", "White"
  ],
  "facialHairType": [
    "Blank", "BeardMedium", "BeardLight", "BeardMagestic", 
    "MoustacheFancy", "MoustacheMagnum"
  ],
  "clotheType": [
    "BlazerShirt", "BlazerSweater", "CollarSweater", 
    "GraphicShirt", "Hoodie", "Overall", 
    "ShirtCrewNeck", "ShirtScoopNeck", "ShirtVNeck"
  ],
  "eyeType": [
    "Close", "Cry", "Default", "Dizzy", "EyeRoll", 
    "Happy", "Hearts", "Side", "Squint", "Surprised", 
    "Wink", "WinkWacky"
  ],
  "eyebrowType": [
    "Angry", "AngryNatural", "Default", "DefaultNatural", 
    "FlatNatural", "RaisedExcited", 
    "RaisedExcitedNatural", "SadConcerned", 
    "SadConcernedNatural", "UnibrowNatural", "UpDown", 
    "UpDownNatural"
  ],
  "mouthType": [
    "Concerned", "Default", "Disbelief", "Eating", 
    "Grimace", "Sad", "ScreamOpen", "Serious", 
    "Smile", "Tongue", "Twinkle", "Vomit"
  ],
  "skinColor": [
    "Tanned", "Yellow", "Pale", "Light", 
    "Brown", "DarkBrown", "Black"
  ],
  "topType": [
    "NoHair", "Eyepatch", "Hat", "Hijab", 
    "Turban", "WinterHat1", "WinterHat2", 
    "WinterHat3", "WinterHat4", "LongHairBigHair", 
    "LongHairBob", "LongHairBun", "LongHairCurly", 
    "LongHairCurvy", "LongHairDreads", "LongHairFrida", 
    "LongHairFro", "LongHairFroBand", "LongHairNotTooLong", 
    "LongHairShavedSides", "LongHairMiaWallace", 
    "LongHairStraight", "LongHairStraight2", 
    "LongHairStraightStrand", "ShortHairDreads01", 
    "ShortHairDreads02", "ShortHairFrizzle", 
    "ShortHairShaggyMullet", "ShortHairShortCurly", 
    "ShortHairShortFlat", "ShortHairShortRound", 
    "ShortHairShortWaved", "ShortHairSides", 
    "ShortHairTheCaesar", "ShortHairTheCaesarSidePart"
  ]
}`

func getRandomElement(elements []string) string {
	return elements[rand.Intn(len(elements))]
}

func generateAvatarURL(data AvatarData) string {
	baseURL := "https://avataaars.io/?"
	params := url.Values{}

	params.Set("accessoriesType", getRandomElement(data.AccessoriesType))
	params.Set("avatarStyle", getRandomElement(data.AvatarStyle))
	params.Set("clotheType", getRandomElement(data.ClotheType))
	params.Set("eyeType", getRandomElement(data.EyeType))
	params.Set("eyebrowType", getRandomElement(data.EyebrowType))
	params.Set("facialHairType", getRandomElement(data.FacialHairType))
	params.Set("hairColor", getRandomElement(data.HairColor))
	params.Set("mouthType", getRandomElement(data.MouthType))
	params.Set("skinColor", getRandomElement(data.SkinColor))
	params.Set("topType", getRandomElement(data.TopType))

	return baseURL + params.Encode()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var data AvatarData
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		http.Error(w, "Failed to marshall JSON", http.StatusInternalServerError)
		return
	}

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	generatedUrl := generateAvatarURL(data)

	// Call the external API that returns an SVG
	resp, err := http.Get(generatedUrl)
	if err != nil {
		http.Error(w, "Failed to fetch SVG", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Set the same Content-Type as the original response
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	// Stream the response from the API to the original response writer
	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
