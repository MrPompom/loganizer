# 🔍 GoLog Analyzer

![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Status](https://img.shields.io/badge/status-active-success.svg)

## 📋 Description

**GoLog Analyzer** est un outil en ligne de commande (CLI) développé en Go, conçu pour aider les administrateurs système à analyser des fichiers de logs provenant de diverses sources (serveurs, applications). L'outil permet de centraliser l'analyse de multiples logs en parallèle et d'en extraire des informations clés, tout en gérant les erreurs de manière robuste.

### 🎯 Objectifs du projet

- **Analyse distribuée** : Traitement concurrent de multiples fichiers de logs
- **Gestion d'erreurs robuste** : Erreurs personnalisées avec gestion fine
- **Interface CLI intuitive** : Utilisation simple avec Cobra
- **Export structuré** : Génération de rapports JSON détaillés
- **Architecture modulaire** : Code organisé en packages logiques

---

## 🚀 Installation

### Prérequis

- **Go 1.19+** installé sur votre système
- Accès en lecture aux fichiers de logs à analyser

### Installation depuis les sources

```bash
# Cloner le repository
git clone https://github.com/votre-username/go_loganizer.git
cd go_loganizer

# Initialiser le module Go
go mod init go_loganizer

# Installer les dépendances
go get github.com/spf13/cobra@latest

# Compiler le projet
go build -o loganalyzer .

# Ou installer globalement
go install .
```

---

## 📖 Utilisation

### Commande principale : `analyze`

La commande `analyze` permet d'analyser une liste de fichiers de logs définis dans un fichier de configuration JSON.

#### Syntaxe

```bash
loganalyzer analyze --config <path_to_config.json> [--output <path_to_output.json>]
```

#### Options

| Flag | Raccourci | Description | Obligatoire |
|------|-----------|-------------|-------------|
| `--config` | `-c` | Chemin vers le fichier de configuration JSON | ✅ |
| `--output` | `-o` | Chemin vers le fichier de sortie pour les résultats | ❌ |

#### Exemples d'utilisation

```bash
# Analyse basique avec configuration
./loganalyzer analyze --config examples/config.json

# Analyse avec export des résultats
./loganalyzer analyze -c examples/config.json -o reports/analysis_report.json

# Affichage de l'aide
./loganalyzer analyze --help
```

---

## 📁 Structure des fichiers

### Fichier de configuration d'entrée

Le fichier de configuration doit être au format JSON et contenir un tableau d'objets représentant les logs à analyser.

**Structure :**

```json
[
  {
    "id": "identifiant-unique",
    "path": "/chemin/vers/le/fichier.log",
    "type": "type-de-log"
  }
]
```

**Exemple (`config.json`) :**

```json
[
  {
    "id": "web-server-1",
    "path": "/var/log/nginx/access.log",
    "type": "nginx-access"
  },
  {
    "id": "app-backend-2", 
    "path": "/var/log/my_app/errors.log",
    "type": "custom-app"
  },
  {
    "id": "system-auth",
    "path": "/var/log/auth.log",
    "type": "system-auth"
  }
]
```

### Fichier de rapport de sortie

Lorsque l'option `--output` est utilisée, un rapport JSON détaillé est généré.

**Structure du rapport :**

```json
[
  {
    "log_id": "identifiant-du-log",
    "file_path": "/chemin/vers/le/fichier",
    "status": "OK|FAILED",
    "message": "Message descriptif",
    "error_details": "Détails de l'erreur si applicable"
  }
]
```

**Exemple de rapport :**

```json
[
  {
    "log_id": "web-server-1",
    "file_path": "/var/log/nginx/access.log",
    "status": "OK",
    "message": "Analyse terminée avec succès.",
    "error_details": ""
  },
  {
    "log_id": "invalid-path",
    "file_path": "/non/existent/log.log", 
    "status": "FAILED",
    "message": "Fichier introuvable.",
    "error_details": "open /non/existent/log.log: no such file or directory"
  }
]
```

---

## ⚙️ Fonctionnalités

### 🔄 Traitement concurrent

- **Goroutines** : Chaque fichier de log est analysé dans une goroutine séparée
- **WaitGroup** : Synchronisation des goroutines pour attendre la fin de tous les traitements
- **Channel sécurisé** : Collecte thread-safe des résultats d'analyse

### 🛡️ Gestion des erreurs personnalisées

Le projet implémente deux types d'erreurs personnalisées :

#### 1. `NonExistingFileError`
- **Description** : Fichier introuvable ou inaccessible
- **Utilisation** : Vérification de l'existence et des permissions d'accès
- **Gestion** : `errors.As()` pour la détection et le traitement spécifique

#### 2. `ParsingError` 
- **Description** : Erreur lors de l'analyse du contenu du fichier
- **Utilisation** : Simulation d'erreurs de parsing (10% de chance)
- **Gestion** : `errors.As()` pour la détection et le traitement spécifique

### 📊 Analyse simulée

- **Temps de traitement** : Simulation aléatoire entre 50ms et 200ms par fichier
- **Taux d'erreur** : 10% de chance d'erreur de parsing simulée
- **Validation** : Vérification de l'existence et de l'accessibilité des fichiers

### 📈 Reporting

- **Console** : Affichage en temps réel des résultats d'analyse
- **Export JSON** : Génération optionnelle d'un rapport détaillé
- **Statuts clairs** : `OK` pour succès, `FAILED` pour échec

---

## 🔧 API des packages internes

### Package `internal/config`

```go
// InputTarget représente une cible d'analyse
type InputTarget struct {
    ID   string `json:"id"`
    Path string `json:"path"`  
    Type string `json:"type"`
}

// LoadTargetsFromFile charge les cibles depuis un fichier JSON
func LoadTargetsFromFile(filePath string) ([]InputTarget, error)
```

### Package `internal/analyzer`

```go
// CheckResult contient le résultat d'une analyse
type CheckResult struct {
    InputTarget InputTarget
    Message     string
    Err         error
}

// ReportEntry représente une entrée du rapport final
type ReportEntry struct {
    LogID        string `json:"log_id"`
    FilePath     string `json:"file_path"`
    Status       string `json:"status"`
    Message      string `json:"message"`
    ErrorDetails string `json:"error_details"`
}

// CheckLog analyse un fichier de log
func CheckLog(target InputTarget) CheckResult

// ConvertToReportEntry convertit un CheckResult en ReportEntry
func ConvertToReportEntry(result CheckResult) ReportEntry
```

### Package `internal/reporter`

```go
// ExportResultToJsonFile exporte les résultats vers un fichier JSON
func ExportResultToJsonFile(filePath string, results []ReportEntry) error
```

---

## 🧪 Tests et exemples

### Fichiers de test fournis

Le projet inclut des fichiers d'exemple dans le dossier `examples/` :

- `config.json` : Configuration type avec plusieurs logs
- `logs/` : Dossier contenant des fichiers de logs de test

### Commandes de test

```bash
# Test avec configuration d'exemple
./loganalyzer analyze -c examples/config.json

# Test avec export
./loganalyzer analyze -c examples/config.json -o test_report.json

# Test avec fichiers inexistants (pour tester la gestion d'erreurs)
./loganalyzer analyze -c examples/config_with_errors.json -o error_report.json
```

---

## 🚨 Gestion des erreurs

### Types d'erreurs gérées

1. **Fichier de configuration invalide**
   - JSON malformé
   - Fichier inexistant
   - Permissions insuffisantes

2. **Fichiers de logs problématiques**
   - Chemin inexistant
   - Permissions d'accès refusées
   - Fichiers corrompus

3. **Erreurs de traitement**
   - Erreurs de parsing simulées
   - Problèmes de mémoire
   - Timeouts

### Messages d'erreur

Les messages d'erreur sont explicites et incluent :
- Le contexte de l'erreur
- L'identifiant du log concerné
- Le chemin du fichier problématique
- Les détails techniques de l'erreur

---

## 🎁 Fonctionnalités bonus (Future)

### Fonctionnalités prévues

1. **Gestion des dossiers d'exportation**
   - Création automatique des répertoires de sortie
   - Support des chemins relatifs et absolus

2. **Horodatage des exports**
   - Nommage automatique avec timestamp
   - Format : `AAMMJJ_report.json`

3. **Commande `add-log`**
   - Ajout interactif de nouvelles configurations
   - Validation en temps réel

4. **Filtrage des résultats**
   - Flag `--status` pour filtrer par statut
   - Support des expressions régulières

---

## 👥 Équipe de développement

| Rôle | Nom | Responsabilités |
|------|-----|----------------|
| **Lead Developer** | [Votre Nom] | Architecture, CLI, Documentation |
| **Backend Developer** | [Nom du coéquipier 1] | Analyzer, Gestion d'erreurs |
| **DevOps** | [Nom du coéquipier 2] | Tests, Configuration, Reporter |

### Contributions

- **Architecture** : Conception modulaire et packages internes
- **Concurrence** : Implémentation des goroutines et synchronisation
- **CLI** : Interface utilisateur avec Cobra
- **Testing** : Scénarios de test et validation
- **Documentation** : README, commentaires de code, exemples

---

## 📚 Ressources et références

### Documentation Go

- [Goroutines et concurrence](https://go.dev/tour/concurrency)
- [Gestion des erreurs](https://go.dev/blog/error-handling-and-go)
- [Package JSON](https://pkg.go.dev/encoding/json)

### Librairies utilisées

- [Cobra CLI](https://github.com/spf13/cobra) - Framework pour applications CLI
- [Go Standard Library](https://pkg.go.dev/std) - Packages standards Go

### Standards de développement

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

---

## 📄 License

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

---

## 🤝 Contributing

Les contributions sont les bienvenues ! Pour contribuer :

1. Fork le projet
2. Créer une branche feature (`git checkout -b feature/AmazingFeature`)
3. Commit vos changements (`git commit -m 'Add AmazingFeature'`)
4. Push vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrir une Pull Request

---

## 📞 Support

Pour toute question ou problème :

- Ouvrir une [issue](https://github.com/votre-username/go_loganizer/issues)
- Contacter l'équipe de développement
- Consulter la [documentation](./docs/)

---

**Made with ❤️ by [Nom de votre équipe]**