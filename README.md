# Programmation répartie 
## Laboratoire 1
### Auteurs 
Tania Nunez, Magali Egger
### Introduction
Ce laboratoire concerne le développement d'une petite application client/server permettant la 
répartition de bénévoles pour l’organisation de manifestations.  
### Utilisation de l'application
Afin d'utiliser notre programme, il faut cloner le repository github puis lancer le serveur et ensuite un/des client(s).

Pour ce faire, ouvrez un terminal à la racine du projet. Lancez le serveur à l'aide de la commande 
```go run ./main/main.go server```.
Lancez le client à l'aide de la commande ```go run ./main/main.go client```

Afin d'arrêter le client, il vous suffit d'entrer la commande ```QUIT``` lors de son exécution.
Afin d'arrêter le server, il faut effectuer un CTRL+C puisqu'il consiste en une boucle infinie à l'écoute de connexions
potentielles.

Si vous changez le port et/ou l'hôte, pensez à effectuer le changement du côté serveur, client et tests.

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

### Fonctionnalités non réalisées
Pour ce travail, il était demandé de charger les utilisateurs et événements depuis un fichier de configuration. Notre
programme possède ces données directement dans la classe dataReaderWriter sous forme de tableau d'utilisateur et
d'événements.
Concernant la possibilité de tester l'accès concurrent aux données manuellement, nous l'avons fait en utilisant des 
breakpoints placés au point critique (par exemple, lorsque quelqu'un crée un event) puis en exécutant un deuxième client
essayant d'effectuer une lecture ou écriture.
