
@Autores: Henry Aliaga / Ismael Espinoza

@Fecha: 12/22/2024

@Descripcion: programa de streaming para musica

@Lenguaje: Golang

@Resumen del programa:

La música ha tenido un impacto significativo en los usuarios, puesto que determina el estado de ánimo, por lo cual se considera la creación de un sistema que permita disfrutar de contenido musical de manera organizada y accesible. Este proyecto surgió como una solución en un ambiente local de streaming musical, diseñada para brindar una experiencia agradable en la reproducción de música.
Adicionalmente la esencia de nuestro sistema es la combinación de algunos lenguajes de programación entre ellos estan Go en el backend con y también HTML, CSS y JavaScript en el frontend, que se encontrara construido en MySQL, con una base de datos relacional.
El sistema se encuentra establecido por 2 roles de usuarios: los administradores, usuarios regulares, que podrán navegar por la biblioteca, y acceder a contenido musical a través de un reproductor.
Se considero que dentro de las limitaciones del entorno local y garantizando que sea un sistema con un funcionamiento fluido, se ha establecido algunos límites en los cuales son, la capacidad de la biblioteca para contener hasta 60 canciones, cada una con un tamaño máximo de 10MB, manteniendo un espacio total de 600MB. Cabe mencionar que el sistema está pensado para funcionar de manera óptima en navegadores modernos, permitiendo que hasta 5 usuarios puedan acceder al sistema de manera simultáneamente.

# Pasos de instalación

1. Clonar/descargar el repositorio https://github.com/SORTERIUS/PROYECTO_STREAMING
2. Ejecutar streaming_music.sql para crear la base de datos y conectarse a MySQL o cualquier gestor de base de datos pero  que sea para BBDD relacionales.
3. Verificar que las canciones estén en Backend/uploads/songs/
4. Eliminar cache e historial de navegador(para evitar conflictos de versiones anteriores en caso de una descarga de un compilado anterior o versionamiento)
5. Iniciar el servidor Go del backend el archivo main.go , con el comando go run main.go en Visual Code