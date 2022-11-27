# Programmation répartie 
## Laboratoire 1
### Auteures 
Tania Nunez, Magali Egger
### Introduction
Ce laboratoire concerne le développement d'une petite application client/server permettant la 
répartition de bénévoles pour l’organisation de manifestations.  
### Utilisation de l'application
Afin d'utiliser notre programme, il faut cloner le repository github puis lancer le serveur et ensuite un/des client(s).

Pour ce faire, ouvrez un terminal à la racine du projet. Lancez le serveur à l'aide de la commande 
```go run ./main/server/mainServer.go server```.
Lancez le client à l'aide de la commande ```go run ./main/client/client.go```.

Afin d'arrêter le client, il vous suffit d'entrer la commande ```QUIT``` lors de son exécution.
Pour arrêter le serveur, il faut effectuer un CTRL+C puisqu'il consiste en une boucle infinie à l'écoute de connexions
potentielles.

Si vous souhaitez lancer un nombre de serveurs souhaité, vous pouvez modifier le main du fichier mainServer.go et
appeler directement la fonction LaunchNServ(n) du package server. Sinon, ce nombre sera celui se trouvant dans le fichier
de configuration. Il y a également la possibilité de lancer un serveur directement à partir de son identifiant en 
appelant la fonction Launch(id, done), notez que le deuxième paramètre est un channel de type bool qui permet de faire
attendre la routine principale jusqu'à ce que le serveur soit lancé et se ferme.

La configuration du programme se trouve dans le fichier ```config.yaml```. Ce fichier permet de modifier les ports utilisés 
ainsi que le nombre de serveurs lancés.

Concernant les tests, il faut vous placer dans le dossier test et exécuter la commande ```go test```.

### Fonctionnalité de l'application
Notre application permet de faire les actions suivantes :
- créer un pool d'utilisateurs, postes et manifestations
- authentification des users
- lister les manifestations
- lister les postes d'une manifestation donnée
- lister les bénévoles ainsi que les postes auquel ils appartiennent d'un événement donné
- Inscrire un utilisateur à un poste dans un événement
- Créer un événement (nécessite d'être authentifié)
- Fermer un événement (nécessite d'être authentifié)
- Quitter l'application

### Explications des commandes 
L'utilisateur peut entrer les commandes suivantes :
- ```ADD``` : permet à un utilisateur de s'inscrire à un poste dans un événement. Cette commande nécessite une authentification, ainsi que de préciser le numéro de la manifestation et le numéro du poste.
- ```LISTM``` : permet de lister les manifestations.
- ```LISTP``` : permet de lister les postes d'une manifestation donnée, en précisant le numéro de la manifestation.
- ```LISTU``` : permet de lister les bénévoles ainsi que les postes auquel ils appartiennent d'un événement donné. Il est nécessaire de préciser le numéro de la manifestation.
- ```CREATE``` : permet de créer un événement. Cette commande nécessite une authentification, ainsi que de préciser le nom de la manifestation, le nom de chaque poste suivie de sa capacité.
- ```CLOSE``` : permet de fermer un événement. Cette commande nécessite une authentification et de préciser le numéro de la manifestation.

### Lancement en mode "debug"
Pour lancer le server en mode debug, on peut modifier le fichier de configuration ```config.json``` en mettant ```"debug": 1```.
Cela permettera de ralentir l'execution du programme afin de pouvoir observer les étapes de l'execution, grâce à un ```time.Sleep()``` de 10 secondes. 
Ensuite il est possible de lancer plusieurs clients en même temps en ouvrant plusieurs terminaux et en effectuant la commande ```go run ./main/client/client.go client```.
En remettant ```"debug": 0```, le programme s'exécutera normalement. 
Pour lancer le serveur en mode trace, il faut lancer la commande ```go run ./main/server/mainServer.go server -trace```. Nous pouvons
faire pareil pour le client avec la commande ```go run ./main/client/client.go client -trace```. 

### Fonctionnement de l'algorithme de Lamport
L'algorithme de Lamport permet de résoudre le problème de l'exclusion mutuelle. Il permet de garantir l'exclusion mutuelle entre les processus server.
Nous avons exploré les différents aspects de l'algorithme de Lamport, notamment la gestion des messages, la gestion de la clock et la gestion de la ressource critique.
Les messages entre les servers ont la forme suivante : ```<type> <clock> <id>```. Le type peut être "req" pour une requête, "rel" pour une release, "ack" pour une acknowledgement, "data" pour une modification
des données (voir en bas du paragraphe) et "ready"
pour indiquer que le serveur est prêt à recevoir des requêtes. La clock est un entier qui représente le temps du serveur. L'id est un entier qui représente l'id du serveur.
Le package processMutex contient les éléments centraux de l'algorithme de Lamport, notamment l'horloge de Lamport implementée dans le fichier lamportClock.go, la gestion de la communication entre
servers dans le fichier network.go et la gestion de l'accès à la ressource critique dans le fichier mutex.go.
L'accès à la section critique se fait via la methode ```AskDataRW()``` du fichier server.go.
Les messages de type data sont envoyés lorsque l'on veut modifier les données. Ils sont envoyés à tous les serveurs commence par le mot "data" suivie du changement à effectuer. Le format du
changement est le même que celle de la commande envoyée par le client.

