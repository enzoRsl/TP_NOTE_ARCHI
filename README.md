# TP Aéroport: Golang / MQTT / Redis

Membres du groupe : Nathan Fouéré, Clément Guibout, Enzo Roussel

Le dépôt du serveur HTTP demandé est [disponible à cette adresse](https://github.com/enzoRsl/airport_web_server).

## Table des matières
- [TP Aéroport: Golang / MQTT / Redis](#tp-aéroport-golang--mqtt--redis)
  - [Table des matières](#table-des-matières)
  - [Instructions d'installation](#instructions-dinstallation)
  - [Instructions de mise en route](#instructions-de-mise-en-route)
  - [Instructions de personnalisation](#instructions-de-personnalisation)
  - [Instructions de build](#instructions-de-build)
  - [Instructions d'exécution](#instructions-dexécution)
  - [Justifications](#justifications)
    - [La génération des données](#la-génération-des-données)
    - [Le fichier de configuration](#le-fichier-de-configuration)
    - [Les identifiants de capteurs](#les-identifiants-de-capteurs)
    - [Pourquoi avoir choisi Redis](#pourquoi-avoir-choisi-redis)
    - [Topic Redis](#topic-redis)
    - [Messages Redis](#messages-redis)
    - [Niveau de QoS](#niveau-de-qos)


## Instructions d'installation

1. Cloner le dépôt

```
$ git clone https://github.com/enzoRsl/TP_NOTE_ARCHI.git
```

2. Se déplacer dans le dépôt local

```
$ cd TP_NOTE_ARCHI
```


## Instructions de mise en route

Commencez par copier le fichier `.env.example` dans un fichier `.env`.
Vous pouvez choisir un nom et un mot de passe différent pour la connexion à MQTT, ainsi que changer l'intervalle (en secondes) de publication des capteurs.

Pour que la base Redis puisse se remplir avec les données envoyées par les pseudo-capteurs, il faut s'assurer d'avoir un serveur Redis actif. Dans notre cadre, nous allons tout héberger sur notre machine.

Lancez votre serveur Redis, puis vérifier qu'il est en route avec la commande ping:

```sh
$ redis-cli -h 127.0.0.1 -p 6379 ping
PONG
```

Si la réponse est `PONG`, votre serveur se comporte correctement. Renseignez maintenant votre host et le port du serveur Redis dans le fichier `.env` dans la variable `REDIS_URI`.

De la même manière, vérifiez que votre broker MQTT est actif, et renseignez son URI dans le fichier `.env`.


## Instructions de personnalisation

Si vous souhaitez personnaliser les options possibles pour les différents capteurs, vous pouvez ajouter des nouveaux aéroports dans la liste déjà présente dans le fichier `config.yml`.

Ce fichier de configuration répertorie tous les aéroports sur lesquels des capteurs peuvent être présents. L'utilité de cette liste est justifiée dans la partie [justifications](#le-fichier-de-config).


## Instructions de build

Windows :

```cmd
$ build.bat
```

Linux / OSX :

```sh
$ ./build.sh
```


## Instructions d'exécution

Pour rendre actif un nouveau capteur (donc par extension remplir la base Redis), vous devez exécuter un ou plusieurs des 3 capteurs existants : `pressure`, `temperature` et `wind`.

Chaque capteur a besoin de deux paramètres à renseigner à son exécution:

1. Un identifiant de capteur, un entier naturel (voir [la justification](#les-identifiants-de-capteurs))
2. Un [code IATA](#https://fr.wikipedia.org/wiki/Code_IATA_des_a%C3%A9roports) présent dans le fichier de configuration `config.yml` (voir [la justification](#le-fichier-de-config))

Si vous souhaitez allumer le capteur de pression numéro 1 sur l'aéroport de Nantes, vous devez suivre la démarche suivante:

```
$ ./build/pressure 1 NTE
```

Suivez le même principe pour les 2 autres types de capteurs.


## Justifications

### La génération des données

Les pseudo-capteurs que nous devions créer pouvaient prendre plusieurs formes : génération aléatoire de données, génération guidée, récupération de données réelles depuis une API...

Pour des raisons de cohérence des donnes, nous avons décidé de récupérer toutes les données météo depuis une API gratuite : l'API [open-meteo](https://open-meteo.com/)

Le problème de cette API est qu'elle n'est mise à jour qu'une fois par heure, alors que la consigne indiquait que les capteurs devaient être beaucoup plus actifs. Nous avons donc décidé de dupliquer les données de l'API, c'est-à-dire qu'une donnée sera répliquée plusieurs fois si l'intervalle de récupération/publication du fichier `.env` est inférieure à une heure.

L'avantage de cette génération de données est qu'elles sont plausibles, puisque les données sont récupérées par de réels capteurs sur l'API.

### Le fichier de configuration

La récupération des données sur une API permet d'obtenir des données réelles, puisque cette API publie des données de capteurs réellement existants. Or, pour obtenir ces données, nous devons lui indiquer à quelle position nous voulons récupérer les informations, grâce à la latitude et longitude.

Cette contrainte nous oblige à lister les aéroports avec leur position (latitude, longitude) dans un fichier de configuration pour que lors de l'exécution de nos capteurs, nous puissions indiquer le code IATA d'un aéroport et que sa position soit automatiquement associée pour que les données demandées à l'API soient cohérentes.


### Les identifiants de capteurs

Nous avons décidé de demander des identifiants de capteurs lors de l'exécution de ceux-ci pour que nous puissions simuler la présence de plusieurs capteurs dans un seul et même aéroport. Par exemple, nous pouvons imaginer que nous ayons deux capteurs de vitesse du vent dans un même aéroport, mais sur deux pistes différentes.

Avec des identifiants de capteurs permet de simuler un capteur numéro `1` sur une première piste, ainsi que `n` autres capteurs sur ce même aéroport.


### Pourquoi avoir choisi Redis

Nous avons décidé de choisir Redis comme base de données. Ce choix vient du fait que nous ne connaissions pas cet outil, et que nous avions envie de découvrir son usage.
De plus, la manière de gérer les données que nous avons utilisées (expliquée plus en détail dans la partie [topic redis](#topic-redis)) peut requérir beaucoup de requêtes à la base pour obtenir certaines données sur une tranche horaires données. Étant donné que Redis est construit sur la RAM, cela nous a permis de ne pas nous soucier de la vitesse d'accès à la base.


### Topic Redis

Les topics Redis sont organisés sous la forme suivante :

```
airport/<code IATA>/datatype/<wind|temperature|pressure>/date/<date au format 2006-01-02>/hour/<heure format 24h>/
```

Exemple :
```
airport/NTE/datatype/wind>/date/2006-01-02/hour/15/
```

Ainsi, nous pouvons récupérer toutes les données qui concernent un aéroport, un type de capteur, une date précise ainsi que l'heure de la journée.


### Messages Redis

Les messages postés dans les topics Redis sont organisés sous la forme suivante :

```
<identifiant du capteur>:<date et heure au format 2006-01-02-15-04-05>:<donnée du capteur dans l'unité métrique internationale>
```

Exemple :
```
1:2006-01-02-15-04-05:14
```

Ainsi, cela nous permet d'avoir différents capteurs avec les données sur le même topic, ce qui est pratique pour agréger les données dans l'API. De plus, nous avons la date directement formatée dans le message.


### Niveau de QoS

Nous avons choisi d'utiliser un niveau de QoS sur le broker MQTT de 2. Nous avons pris cette décision car de la manière que nous avons organisé notre Redis, il serait possible d'avoir des doublons de mesure de capteur, puisque les clés ne sont pas uniques.

Si le niveau de QoS était de 0, nous n'aurions pas de garantie que le message a bien été délivré. Nous avons jugé des capteurs d'aéroport trop sensibles pour risquer de perdre des informations.

Si le niveau de QoS était de 1, le broker MQTT aurait pu recevoir plusieurs fois le même message, et donc fausser les moyennes et les données stockées.

Nous n'avions donc d'autre choix que d'avoir un niveau de QoS de 2 avec la méthode que nous avons décidé d'utiliser.
