# Programmation répartie 
## Laboratoire 1
### Auteurs 
Tania Nunez, Magali Egger
### Introduction
Ce laboratoire concerne le développement d'une petite application client/server permettant la 
répartition de bénévoles pour l’organisation de manifestations.  
### Utilisation de l'application
Afin d'utiliser notre programme, il suffit simplement de commencer par lancer le server depuis un terminal,
en allant à l'endroit ou se trouve le fichier server.go, et en écrivant la commande : go run server.go.
Ensuite il est possible de lancer autant de client que souhaiter en se rendant à l'endroit ou se trouve le fichier 
client.go, et en écrivant la commande : go run client.go.

Ouvrez un terminal à la racine du projet. Lancez le serveur à l'aide de la commande ``go run ./main/server/server.go``.
Lancez le client à l'aide de la commande ``go run ./main/client/client.go``.

Afin d'arrêter le client, il vous suffit d'entrer la commande ``QUIT`` lors de son exécution.
Afin d'arrêter le server, il faut effectuer un CTRL+C puisqu'il consiste en une boucle infinie à l'écoute de connexions
potentielles.

Si vous changez le port et/ou l'hôte, pensez à effectuer le changement du côté serveur et client.

### Fonctionnalité de l'application
Notre application permet de faire les actions suivantes :
- créer un pool d'utilisateurs, postes et manifestations
- authentification des users
- lister les manifestations
- lister les postes d'une manifestation donnée
- lister les bénévoles ainsi que les postes auquel ils appartiennent
- Inscrire un utilisateur à un poste dans un événement

### Fonctionnalités non réalisées
Pour ce travail, il était demandé de charger les utilisateurs et événements depuis un fichier de configuration. Notre
programme possède ces données directement dans la classe dataReaderWriter sous forme de tableau d'utilisateur et
d'événements.
Concernant la possibilité de tester l'accès concurrent aux données manuellement, nous l'avons fait en utilisant des 
breakpoints placés au point critique (par exemple, lorsque quelqu'un crée un event) puis en exécutant un deuxième client
essayant d'effectuer une lecture ou écriture.
