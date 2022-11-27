
### Probleme a regler du premier laboratoire
1. Test -> modifier le readme pour mode debug
2. ~~Il aurait été mieux d’avoir 2 « main » complètement séparés, un qui s’occupe exclusivement du serveur, et l’autre du client. C’est plus judicieux et plus propre (-1).~~
3. ~~Il y a un souci lorsque qu’un client vient d’arriver et s’en va avec la commande QUIT, ou Control-C. Dans ce cas-là, le serveur s’arrête complètement. Ceci ne devrait évidemment pas arriver (-1).~~
4. Le serveur ne devrait pas s’occuper de produire ce que le client va afficher mais envoyer uniquement au client les données utiles à un affichage par ce dernier (-1)
5. ~~Les commandes du programme ne sont pas intuitives (quand on voit ADD, on sait pas ce que ça veut dire…). Expliquez vos commandes dans le README svp (-1).~~
6. ~~fichier de configuration JSON pour précharger des données existantes~~
7. ~~Lorsqu’un bénévole aimerait s’inscrire dans un poste (de l’une des manifestations) dont le nombre maximum de bénévoles a été atteint, alors l’inscription doit être refusée.~~
8. ~~Lorsqu’un événement a été fermé, votre serveur accepte toujours les inscriptions d’un bénévole alors que ça ne doit pas être le cas~~
9. ~~La commande qui affiche les manifestations (LISTM) doit également afficher les manifestations qui ont été clôturées aux inscriptions par leur organisateur. Ce n’est pas le cas dans votre implémentation~~
