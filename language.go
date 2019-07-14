package gotagger

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var supported = map[string][]string{
	"es": {"a", "actualmente", "acuerdo", "adelante", "ademas", "además",
		"adrede", "afirmó", "agregó", "ahi", "ahora", "ahí", "al", "algo",
		"alguna", "algunas", "alguno", "algunos", "algún", "alli", "allí",
		"alrededor", "ambos", "ampleamos", "antano", "antaño", "ante",
		"anterior", "antes", "apenas", "aproximadamente", "aquel", "aquella",
		"aquellas", "aquello", "aquellos", "aqui", "aquél", "aquélla",
		"aquéllas", "aquéllos", "aquí", "arriba", "arribaabajo", "aseguró",
		"asi", "así", "atras", "aun", "aunque", "ayer", "añadió", "aún", "b",
		"bajo", "bastante", "bien", "breve", "buen", "buena", "buenas", "bueno",
		"buenos", "c", "cada", "casi", "cerca", "cierta", "ciertas", "cierto",
		"ciertos", "cinco", "claro", "comentó", "como", "con", "conmigo",
		"conocer", "conseguimos", "conseguir", "considera", "consideró",
		"consigo", "consigue", "consiguen", "consigues", "contigo", "contra",
		"cosas", "creo", "cual", "cuales", "cualquier", "cuando", "cuanta",
		"cuantas", "cuanto", "cuantos", "cuatro", "cuenta", "cuál", "cuáles",
		"cuándo", "cuánta", "cuántas", "cuánto", "cuántos", "cómo", "d", "da",
		"dado", "dan", "dar", "de", "debajo", "debe", "deben", "debido",
		"decir", "dejó", "del", "delante", "demasiado", "demás", "dentro",
		"deprisa", "desde", "despacio", "despues", "después", "detras",
		"detrás", "dia", "dias", "dice", "dicen", "dicho", "dieron",
		"diferente", "diferentes", "dijeron", "dijo", "dio", "donde", "dos",
		"durante", "día", "días", "dónde", "e", "ejemplo", "el", "ella",
		"ellas", "ello", "ellos", "embargo", "empleais", "emplean", "emplear",
		"empleas", "empleo", "en", "encima", "encuentra", "enfrente",
		"enseguida", "entonces", "entre", "era", "erais", "eramos", "eran",
		"eras", "eres", "es", "esa", "esas", "ese", "eso", "esos", "esta",
		"estaba", "estabais", "estaban", "estabas", "estad", "estada",
		"estadas", "estado", "estados", "estais", "estamos", "estan", "estando",
		"estar", "estaremos", "estará", "estarán", "estarás", "estaré",
		"estaréis", "estaría", "estaríais", "estaríamos", "estarían",
		"estarías", "estas", "este", "estemos", "esto", "estos", "estoy",
		"estuve", "estuviera", "estuvierais", "estuvieran", "estuvieras",
		"estuvieron", "estuviese", "estuvieseis", "estuviesen", "estuvieses",
		"estuvimos", "estuviste", "estuvisteis", "estuviéramos", "estuviésemos",
		"estuvo", "está", "estábamos", "estáis", "están", "estás", "esté",
		"estéis", "estén", "estés", "ex", "excepto", "existe", "existen",
		"explicó", "expresó", "f", "fin", "final", "fue", "fuera", "fuerais",
		"fueran", "fueras", "fueron", "fuese", "fueseis", "fuesen", "fueses",
		"fui", "fuimos", "fuiste", "fuisteis", "fuéramos", "fuésemos", "g",
		"general", "gran", "grandes", "gueno", "h", "ha", "haber", "habia",
		"habida", "habidas", "habido", "habidos", "habiendo", "habla", "hablan",
		"habremos", "habrá", "habrán", "habrás", "habré", "habréis", "habría",
		"habríais", "habríamos", "habrían", "habrías", "habéis", "había",
		"habíais", "habíamos", "habían", "habías", "hace", "haceis", "hacemos",
		"hacen", "hacer", "hacerlo", "haces", "hacia", "haciendo", "hago",
		"han", "has", "hasta", "hay", "haya", "hayamos", "hayan", "hayas",
		"hayáis", "he", "hecho", "hemos", "hicieron", "hizo", "horas", "hoy",
		"hube", "hubiera", "hubierais", "hubieran", "hubieras", "hubieron",
		"hubiese", "hubieseis", "hubiesen", "hubieses", "hubimos", "hubiste",
		"hubisteis", "hubiéramos", "hubiésemos", "hubo", "i", "igual",
		"incluso", "indicó", "informo", "informó", "intenta", "intentais",
		"intentamos", "intentan", "intentar", "intentas", "intento", "ir", "j",
		"junto", "k", "l", "la", "lado", "largo", "las", "le", "lejos", "les",
		"llegó", "lleva", "llevar", "lo", "los", "luego", "lugar", "m", "mal",
		"manera", "manifestó", "mas", "mayor", "me", "mediante", "medio",
		"mejor", "mencionó", "menos", "menudo", "mi", "mia", "mias", "mientras",
		"mio", "mios", "mis", "misma", "mismas", "mismo", "mismos", "modo",
		"momento", "mucha", "muchas", "mucho", "muchos", "muy", "más", "mí",
		"mía", "mías", "mío", "míos", "n", "nada", "nadie", "ni", "ninguna",
		"ningunas", "ninguno", "ningunos", "ningún", "no", "nos", "nosotras",
		"nosotros", "nuestra", "nuestras", "nuestro", "nuestros", "nueva",
		"nuevas", "nuevo", "nuevos", "nunca", "o", "ocho", "os", "otra",
		"otras", "otro", "otros", "p", "pais", "para", "parece", "parte",
		"partir", "pasada", "pasado", "paìs", "peor", "pero", "pesar", "poca",
		"pocas", "poco", "pocos", "podeis", "podemos", "poder", "podria",
		"podriais", "podriamos", "podrian", "podrias", "podrá", "podrán",
		"podría", "podrían", "poner", "por", "por qué", "porque", "posible",
		"primer", "primera", "primero", "primeros", "principalmente", "pronto",
		"propia", "propias", "propio", "propios", "proximo", "próximo",
		"próximos", "pudo", "pueda", "puede", "pueden", "puedo", "pues", "q",
		"qeu", "que", "quedó", "queremos", "quien", "quienes", "quiere",
		"quiza", "quizas", "quizá", "quizás", "quién", "quiénes", "qué", "r",
		"raras", "realizado", "realizar", "realizó", "repente", "respecto", "s",
		"sabe", "sabeis", "sabemos", "saben", "saber", "sabes", "sal", "salvo",
		"se", "sea", "seamos", "sean", "seas", "segun", "segunda", "segundo",
		"según", "seis", "ser", "sera", "seremos", "será", "serán", "serás",
		"seré", "seréis", "sería", "seríais", "seríamos", "serían", "serías",
		"seáis", "señaló", "si", "sido", "siempre", "siendo", "siete", "sigue",
		"siguiente", "sin", "sino", "sobre", "sois", "sola", "solamente",
		"solas", "solo", "solos", "somos", "son", "soy", "soyos", "su",
		"supuesto", "sus", "suya", "suyas", "suyo", "suyos", "sé", "sí", "sólo",
		"t", "tal", "tambien", "también", "tampoco", "tan", "tanto", "tarde",
		"te", "temprano", "tendremos", "tendrá", "tendrán", "tendrás", "tendré",
		"tendréis", "tendría", "tendríais", "tendríamos", "tendrían",
		"tendrías", "tened", "teneis", "tenemos", "tener", "tenga", "tengamos",
		"tengan", "tengas", "tengo", "tengáis", "tenida", "tenidas", "tenido",
		"tenidos", "teniendo", "tenéis", "tenía", "teníais", "teníamos",
		"tenían", "tenías", "tercera", "ti", "tiempo", "tiene", "tienen",
		"tienes", "toda", "todas", "todavia", "todavía", "todo", "todos",
		"total", "trabaja", "trabajais", "trabajamos", "trabajan", "trabajar",
		"trabajas", "trabajo", "tras", "trata", "través", "tres", "tu", "tus",
		"tuve", "tuviera", "tuvierais", "tuvieran", "tuvieras", "tuvieron",
		"tuviese", "tuvieseis", "tuviesen", "tuvieses", "tuvimos", "tuviste",
		"tuvisteis", "tuviéramos", "tuviésemos", "tuvo", "tuya", "tuyas",
		"tuyo", "tuyos", "tú", "u", "ultimo", "un", "una", "unas", "uno",
		"unos", "usa", "usais", "usamos", "usan", "usar", "usas", "uso",
		"usted", "ustedes", "v", "va", "vais", "valor", "vamos", "van",
		"varias", "varios", "vaya", "veces", "ver", "verdad", "verdadera",
		"verdadero", "vez", "vosotras", "vosotros", "voy", "vuestra",
		"vuestras", "vuestro", "vuestros", "w", "x", "y", "ya", "yo", "z", "él",
		"éramos", "ésa", "ésas", "ése", "ésos", "ésta", "éstas", "éste",
		"éstos", "última", "últimas", "último", "últimos",
	},
	"en": {"a", "about", "above", "after", "again", "against", "all", "am",
		"an", "and", "any", "are", "aren't", "as", "at", "be", "because",
		"been", "before", "being", "below", "between", "both", "but", "by",
		"can't", "cannot", "could", "couldn't", "did", "didn't", "do", "does",
		"doesn't", "doing", "don't", "down", "during", "each", "few", "for",
		"from", "further", "had", "hadn't", "has", "hasn't", "have", "haven't",
		"having", "he", "he'd", "he'll", "he's", "her", "here", "here's",
		"hers", "herself", "him", "himself", "his", "how", "how's", "i", "i'd",
		"i'll", "i'm", "i've", "if", "in", "into", "is", "isn't", "it", "it's",
		"its", "itself", "let's", "me", "more", "most", "mustn't", "my",
		"myself", "no", "nor", "not", "of", "off", "on", "once", "only", "or",
		"other", "ought", "our", "ours", "ourselves", "out", "over", "own",
		"same", "shan't", "she", "she'd", "she'll", "she's", "should",
		"shouldn't", "so", "some", "such", "than", "that", "that's", "the",
		"their", "theirs", "them", "themselves", "then", "there", "there's",
		"these", "they", "they'd", "they'll", "they're", "they've", "this",
		"those", "through", "to", "too", "under", "until", "up", "very", "was",
		"wasn't", "we", "we'd", "we'll", "we're", "we've", "were", "weren't",
		"what", "what's", "when", "when's", "where", "where's", "which",
		"while", "who", "who's", "whom", "why", "why's", "with", "won't",
		"would", "wouldn't", "you", "you'd", "you'll", "you're", "you've",
		"your", "yours", "yourself", "yourselves",
	},
}

// Struct to define language object with its code and stopwords
type language struct {
	code      string
	stopwords []string
}

// loadLanguage function loads language checking 'STOPWORDS' environment
// variable path and loading list from local storage if exists or assigns
// default list. Receives language code. Return language struct or error.
func loadLanguage(code string) (l language, e error) {
	l = language{code, []string{}}

	if env := os.Getenv("STOPWORDS"); env != "" {
		var f string = filepath.Join(env, code)

		var fd *os.File
		if fd, e = os.Open(f); e != nil {
			return l, e
		}
		defer fd.Close()

		var s *bufio.Scanner = bufio.NewScanner(fd)
		s.Split(bufio.ScanLines)

		for s.Scan() {
			var w string = s.Text()
			if len(w) > 0 {
				l.stopwords = append(l.stopwords, w)
			}
		}

	} else {
		var ok bool
		if l.stopwords, ok = supported[code]; !ok {
			return l, errors.New("language not supported")
		}
	}

	return l, e
}

func (l language) isStopword(s string) bool {
	var is = false
	var _s = strings.ToLower(s)
	for _, stw := range l.stopwords {
		is = is || stw == _s
	}

	return is
}
