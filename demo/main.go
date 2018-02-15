package main

import (
	"github.com/lucasmenendez/gotokenizer"
	"github.com/lucasmenendez/gotagger"
	"fmt"
)

func main() {
	text := `
¿Por qué tiene que preocuparle a alguien, sobre todo a aquellos ciudadanos que no pertenecen a un sindicato (que son la mayoría), que la afiliación sindical se sitúe en unos de los niveles más bajos de su historia, y que es posible que continúe descendiendo? A veces hay que recordar lo obvio: porque los sindicatos son la institución que mejor garantiza que los asalariados tengan una voz fuerte que los represente tanto en el mercado como en la democracia (a través de la defensa de los derechos sociales). Cuando los sindicatos son fuertes están en condiciones de garantizar que se paguen salarios más justos, que se los tenga en cuenta en las decisiones que toman las empresas o la Administración, o que reciban la formación que precisan para ascender en la escala de la movilidad social. También pueden facilitar la creación de empleo. Por último, pero no menos importante, los sindicatos pueden fomentar la participación política y ayudar a los asalariados a conseguir políticas públicas como el salario mínimo y los instrumentos de protección social del Estado de Bienestar.

Hasta aquí la teoría. ¿Qué está ocurriendo en nuestro país, y en general en el mundo occidental en relación a esta parte tan esencial de la sociedad civil? Un día –es una utopía factible– se corregirá la gangrena de la corrupción y encontraremos la fórmula de la conllevancia entre Cataluña y el resto de España. Entonces, aparecerán en primera fila, con toda su crudeza, los verdaderos problemas estructurales. El más grave de ellos es que nos hemos convertido en un país de precarios, aquejados de una inseguridad creciente, en todas sus cohortes de edad: jóvenes, mujeres, mayores de 45 años, jubilados, etcétera. Ello afecta no sólo al bienestar –ya sería suficiente para combatirla– sino a la calidad de la misma democracia. No se puede ser libre con un temor creciente.

Este panorama se debe a una serie de factores, internos y exteriores. Entre los primeros están las continuas reformas laborales, todas ellas en la misma dirección. La más dañina fue la aprobada por el PP en el año 2012, fundamentalmente porque su objetivo (conseguido) era desequilibrar el poder de cada parte en el seno de la empresa: menor poder para los representantes de los trabajadores, mayor para el empresario. Ello se hizo transformando el marco regulador de los derechos y obligaciones que se desprenden del mercado de trabajo: las condiciones de contratación, la materias relacionadas con la fijación del salario y de las condiciones de trabajo, las extinciones de contratos, etcétera; y también afectando a la dimensión colectiva de las relaciones de trabajo, mediante la modificación del régimen jurídico de la negociación colectiva (con efecto sobre la duración de los convenios, la eficacia de las condiciones acordadas y la estructura de las mismas).

Ello se hizo aprovechando un contexto de extremas dificultades: por una parte, la crisis económica, con sus secuelas de empobrecimiento, desigualdad, precarización y reducción de la protección social; por otra parte, la explosión de un mundo digital y la robotización, que ha fragmentado en mil pedazos el mercado laboral, eliminando la sensación de pertenencia a una comunidad de millones de trabajadores. Como consecuencia de estas dos circunstancias, en estos momentos se están destruyendo más empleos de los que se crean a nivel agregado mundial.

Los sindicatos han tenido que responder a esta secuencia de crisis económica más revolución tecnológica con escasos medios: son instituciones muy representativas pero con escaso nivel de afiliación. Y, como el resto de la sociedad, han ido detrás de las circunstancias (de ahí, buena parte de las acusaciones de que no hacen nada, de que no han estado en la calle con los indignados, de que sólo representan a los trabajadores con empleo y se han olvidado de los jóvenes y de los que laboran en las plataformas digitales en circunstancias muy distintas, que forman parte de una especie de lumpen-proletariado por sus condiciones de vida y de trabajo, etcétera).

En cierta ocasión, y para otro terreno de juego, Manuel Vázquez Montalbán escribió sobre la “correlación de debilidades”. ¿Qué sentido tiene que ante esta gran transformación que se lleva por delante mucho de lo construido después de la Segunda Guerra Mundial en materia de derechos sociales y económicos y de protección social, siga habiendo en España dos estructuras sindicales centrales de clase (CCOO y UGT) –que se enfrentan a una sola patronal, la CEOE– y no se hayan iniciado, hace tiempo ya, las reflexiones y los estudios para la unidad orgánica del movimiento sindical de clase, que lo fortalezca para dar respuesta a lo que tiene encima y lo que llega? Ahora que los sindicatos ya no son correas de transmisión de los partidos políticos como hasta hace poco tiempo (CCOO, del Partido Comunista; UGT, de los socialistas), deja de existir esa dificultad política.

Seguramente las direcciones y las burocracias de los sindicatos conozcan los pros y los contras de esta propuesta de unidad orgánica sindical, que de vez en cuanto emerge y se sumerge como una serpiente de verano, pero el resto de la ciudadanía querría ver un debate abierto sobre las mejores fórmulas para conseguir que el sindicalismo sea más fuerte y esté mejor preparado, aumente su capacidad de afiliación y de atractivo entre los jóvenes, las mujeres, los cuadros, los autónomos, los becarios, los pensionistas, y atiendan a ese mundo creciente de asalariados sin centro de trabajo, esa agregación de colectivos unidos por su inseguridad.

¿Por qué no se abre ya este debate? ¿Por qué la unidad sindical es un silencio social más? ¿Por qué no se avanza en un sindicalismo europeo, internacional de verdad, capaz de influir en el nivel de las causas y no sólo en el de las consecuencias? ¿Quién está interesado en la división? ¿Cambiaría con la unidad la correlación de fuerzas, tan desigual ahora, tan asimétrica, en los diferentes segmentos del sistema económico? Nosotros pensamos que sí, pero mientras tanto estamos abiertos al debate.
`
	sentences := gotokenizer.Sentences(text)

	var words [][]string
	for _, s := range sentences {
		words = append(words, gotokenizer.Words(s))
	}

	if tags, err := gotagger.GetTags(words, "es", 10); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%q\n", tags)
	}
}
