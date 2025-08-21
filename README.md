# 🚀 Git Assistant Intelligent

Le **Git Assistant Intelligent** est un outil en ligne de commande (CLI) écrit en Go qui simplifie et accélère l'utilisation de Git pour les développeurs. Il propose une interface interactive et des commandes automatisées pour les tâches Git les plus courantes, rendant la gestion des dépôts plus intuitive et moins sujette aux erreurs.

![Texte alternatif](/screen.png)

## ✨ Fonctionnalités

  * **Statut Intelligent** : Affiche un résumé clair du statut du dépôt (branche actuelle, commits, fichiers, etc.) et des changements en cours.
  * **Commit Rapide** : Permet d'ajouter et de commiter les changements en une seule étape, avec une sélection de messages de commit prédéfinis.
  * **Gestion des Branches** : Créez, supprimez, changez ou fusionnez des branches avec des commandes simplifiées, adaptées à des flux de travail de développement (ex: `feature/`, `bugfix/`).
  * **Historique Interactif** : Explorez l'historique des commits, visualisez les détails des commits, effectuez des resets ou créez de nouvelles branches à partir de n'importe quel commit.
  * **Analyse de Projet** : Obtenez des informations utiles sur votre dépôt, telles que le nombre de commits, les types de fichiers, et l'activité récente.
  * **Navigation Facile** : Changez de répertoire de travail directement depuis l'application.
  * **Interface Intuitive** : Une interface conviviale avec des menus clairs et des couleurs pour une meilleure lisibilité.

## ⚙️ Comment l'utiliser

### Prérequis

  * [Go](https://golang.org/dl/) installé (version 1.16 ou supérieure recommandée)
  * [Git](https://git-scm.com/downloads) installé et configuré sur votre machine

### Exécution

Pour exécuter l'assistant, suivez ces étapes :

1.  Clonez ce dépôt (ou téléchargez le fichier `GitCtrl.go`).
2.  Ouvrez un terminal et naviguez jusqu'au dossier contenant `GitCtrl.go`.
3.  Lancez l'application avec la commande :
    ```bash
    go run GitCtrl.go
    ```
4.  L'assistant vous demandera de définir votre répertoire de travail. Entrez le chemin d'un dépôt Git existant ou d'un nouveau dossier pour l'initialiser.

### Guide des commandes

Une fois l'assistant lancé, vous serez accueilli par un menu principal.

  * **1. ⚡ Commit rapide** : Ajoute tous les fichiers modifiés et non suivis et les commite.
  * **2. 🌿 Gestion intelligente des branches** : Ouvre un sous-menu pour les opérations de branche.
  * **3. 📜 Historique interactif** : Affiche le log des 15 derniers commits et propose des actions comme le `diff` ou le `reset`.
  * **4. 📊 Analyse du projet** : Donne des statistiques sur votre dépôt.
  * **5. 📁 Changer de répertoire** : Modifie le répertoire de travail de l'application.
  * **6. 🔧 Initialiser Git** : Initialise un nouveau dépôt Git dans le répertoire actuel.
  * **0. ❌ Quitter** : Ferme l'application.

## 🤝 Contribution

Les contributions sont les bienvenues \! Si vous avez des suggestions, des rapports de bugs ou des idées de nouvelles fonctionnalités, n'hésitez pas à ouvrir une *issue* ou à soumettre une *pull request*.

## 📄 Licence

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.
