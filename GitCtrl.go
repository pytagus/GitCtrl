package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Couleurs ANSI
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorCyan   = "\033[36m" // Bleu clair
	ColorBold   = "\033[1m"
)

// Fonctions utilitaires pour les couleurs
func colorize(color, text string) string {
	return color + text + ColorReset
}

func green(text string) string {
	return colorize(ColorGreen, text)
}

func red(text string) string {
	return colorize(ColorRed, text)
}

func cyan(text string) string {
	return colorize(ColorCyan, text)
}

func bold(text string) string {
	return ColorBold + text + ColorReset
}

type GitAssistant struct {
	workingDir   string
	quickCommits []string
	lastActions  []string
}

func NewGitAssistant() *GitAssistant {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Erreur lors de la récupération du répertoire courant:", err)
	}
	return &GitAssistant{
		workingDir: wd,
		quickCommits: []string{
			"🚀 Mise à jour rapide",
			"🐛 Correction de bug",
			"✨ Nouvelle fonctionnalité",
			"📝 Documentation",
			"♻️ Refactoring",
			"🎨 Améliorations UI",
			"⚡ Performance",
			"🔧 Configuration",
		},
		lastActions: make([]string, 0),
	}
}

func (ga *GitAssistant) runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = ga.workingDir
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (ga *GitAssistant) addToHistory(action string) {
	ga.lastActions = append(ga.lastActions, fmt.Sprintf("[%s] %s", time.Now().Format("15:04"), action))
	if len(ga.lastActions) > 10 {
		ga.lastActions = ga.lastActions[1:]
	}
}

func (ga *GitAssistant) isGitRepo() bool {
	_, err := os.Stat(filepath.Join(ga.workingDir, ".git"))
	return err == nil
}

func (ga *GitAssistant) getCurrentBranch() string {
	output, err := ga.runCommand("git", "branch", "--show-current")
	if err != nil {
		return "main"
	}
	return strings.TrimSpace(output)
}

func (ga *GitAssistant) getRepoStats() (int, int, int) {
	// Commits count
	commitsOutput, _ := ga.runCommand("git", "rev-list", "--count", "HEAD")
	commits, _ := strconv.Atoi(strings.TrimSpace(commitsOutput))
	
	// Files count
	filesOutput, _ := ga.runCommand("git", "ls-files")
	files := len(strings.Split(strings.TrimSpace(filesOutput), "\n"))
	if filesOutput == "" {
		files = 0
	}
	
	// Branches count
	branchesOutput, _ := ga.runCommand("git", "branch")
	branches := len(strings.Split(strings.TrimSpace(branchesOutput), "\n"))
	if branchesOutput == "" {
		branches = 0
	}
	
	return commits, files, branches
}

func (ga *GitAssistant) smartStatus() error {
	fmt.Println("📊 === STATUT INTELLIGENT ===")
	
	// Infos générales
	currentBranch := ga.getCurrentBranch()
	commits, files, branches := ga.getRepoStats()
	
	fmt.Printf("🌿 Branche actuelle: %s\n", currentBranch)
	fmt.Printf("📁 Projet: %s\n", filepath.Base(ga.workingDir))
	fmt.Printf("📊 %d commits | %d fichiers | %d branches\n\n", commits, files, branches)
	
	// Changements en cours
	status, err := ga.runCommand("git", "status", "--porcelain")
	if err != nil {
		return err
	}
	
	if strings.TrimSpace(status) == "" {
		fmt.Println("✅ Aucun changement - Dépôt propre")
		
		// Dernier commit
		lastCommit, err := ga.runCommand("git", "log", "-1", "--pretty=format:%h - %s (%cr)")
		if err == nil && lastCommit != "" {
			fmt.Printf("📝 Dernier commit: %s\n", lastCommit)
		}
	} else {
		ga.analyzeChanges(status)
	}
	
	return nil
}

func (ga *GitAssistant) analyzeChanges(status string) {
	lines := strings.Split(strings.TrimSpace(status), "\n")
	
	var added, modified, deleted, untracked []string
	
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		
		statusCode := line[:2]
		filename := line[3:]
		
		switch statusCode[0] {
		case 'A':
			added = append(added, filename)
		case 'M':
			modified = append(modified, filename)
		case 'D':
			deleted = append(deleted, filename)
		case '?':
			untracked = append(untracked, filename)
		}
	}
	
	fmt.Println("📝 Changements détectés:")
	
	if len(added) > 0 {
		fmt.Printf("  ✅ Ajoutés (%d): %s\n", len(added), strings.Join(added, ", "))
	}
	if len(modified) > 0 {
		fmt.Printf("  ✏️ Modifiés (%d): %s\n", len(modified), strings.Join(modified, ", "))
	}
	if len(deleted) > 0 {
		fmt.Printf("  🗑️ Supprimés (%d): %s\n", len(deleted), strings.Join(deleted, ", "))
	}
	if len(untracked) > 0 {
		fmt.Printf("  📂 Non suivis (%d): %s\n", len(untracked), strings.Join(untracked, ", "))
	}
	
	// Suggestions intelligentes
	fmt.Println("\n💡 Suggestions:")
	if len(untracked) > 0 || len(modified) > 0 || len(added) > 0 {
		fmt.Println("  → Utilisez 'Sync rapide' pour ajouter et commiter automatiquement")
	}
	if len(deleted) > 0 {
		fmt.Println("  → Des fichiers ont été supprimés - vérifiez que c'est intentionnel")
	}
}

func (ga *GitAssistant) quickCommit() error {
	// Vérifier s'il y a des changements
	status, err := ga.getStatus()
	if err != nil {
		return err
	}
	
	if strings.TrimSpace(status) == "" {
		fmt.Println("ℹ️ Aucun changement à commiter")
		return nil
	}
	
	fmt.Printf("🚀 === %s ===\n", bold("COMMIT RAPIDE"))
	fmt.Println("Messages prédéfinis:")
	
	for i, msg := range ga.quickCommits {
		fmt.Printf("%d. %s\n", i+1, green(msg))
	}
	fmt.Printf("%d. %s\n", len(ga.quickCommits)+1, cyan("💬 Message personnalisé"))
	
	fmt.Print(cyan("\nChoisissez (1-9): "))
	choice := ga.getUserInput()
	
	var message string
	if choice == strconv.Itoa(len(ga.quickCommits)+1) {
		fmt.Print(cyan("💬 Votre message: "))
		message = ga.getUserInput()
	} else {
		idx, err := strconv.Atoi(choice)
		if err != nil || idx < 1 || idx > len(ga.quickCommits) {
			fmt.Println(red("❌ Choix invalide, utilisation du message par défaut"))
			message = ga.quickCommits[0]
		} else {
			message = ga.quickCommits[idx-1]
		}
	}
	
	// Auto-add et commit
	if err := ga.addAll(); err != nil {
		return err
	}
	
	if err := ga.commit(message); err != nil {
		return err
	}
	
	ga.addToHistory(fmt.Sprintf("Commit rapide: %s", message))
	return nil
}

func (ga *GitAssistant) intelligentBranching() error {
	fmt.Printf("🌿 === %s ===\n", bold("GESTION INTELLIGENTE DES BRANCHES"))
	
	// Lister toutes les branches avec infos
	branchesOutput, err := ga.runCommand("git", "branch", "-v")
	if err != nil {
		return err
	}
	
	fmt.Printf("%s:\n", cyan("Branches existantes"))
	lines := strings.Split(strings.TrimSpace(branchesOutput), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			if strings.HasPrefix(line, "*") {
				// Branche active en vert
				fmt.Printf("%s\n", green(line))
			} else {
				fmt.Printf("%s\n", line)
			}
		}
	}
	
	fmt.Printf("\n%s:\n", cyan("Actions disponibles"))
	fmt.Println("1. 🌱 Créer branche de fonctionnalité")
	fmt.Println("2. 🐛 Créer branche de correction")
	fmt.Println("3. 🔄 Changer de branche")
	fmt.Println("4. 🗑️ Supprimer une branche")
	fmt.Println("5. 🔀 Fusionner une branche")
	fmt.Print(cyan("\nChoisissez (1-5): "))
	
	choice := ga.getUserInput()
	
	switch choice {
	case "1":
		return ga.createFeatureBranch()
	case "2":
		return ga.createBugfixBranch()
	case "3":
		return ga.switchBranch("")
	case "4":
		return ga.deleteBranch()
	case "5":
		return ga.mergeBranch()
	default:
		fmt.Println(red("❌ Choix invalide"))
	}
	
	return nil
}

func (ga *GitAssistant) createFeatureBranch() error {
	fmt.Print("✨ Nom de la fonctionnalité: ")
	feature := ga.getUserInput()
	if feature == "" {
		return fmt.Errorf("nom requis")
	}
	
	// Nettoyer le nom
	branchName := "feature/" + strings.ReplaceAll(strings.ToLower(feature), " ", "-")
	
	_, err := ga.runCommand("git", "checkout", "-b", branchName)
	if err != nil {
		return err
	}
	
	fmt.Printf("✅ Branche '%s' créée et activée!\n", branchName)
	ga.addToHistory(fmt.Sprintf("Branche créée: %s", branchName))
	return nil
}

func (ga *GitAssistant) createBugfixBranch() error {
	fmt.Print("🐛 Description du bug: ")
	bug := ga.getUserInput()
	if bug == "" {
		return fmt.Errorf("description requise")
	}
	
	branchName := "bugfix/" + strings.ReplaceAll(strings.ToLower(bug), " ", "-")
	
	_, err := ga.runCommand("git", "checkout", "-b", branchName)
	if err != nil {
		return err
	}
	
	fmt.Printf("✅ Branche '%s' créée et activée!\n", branchName)
	ga.addToHistory(fmt.Sprintf("Branche de correction créée: %s", branchName))
	return nil
}

func (ga *GitAssistant) deleteBranch() error {
	fmt.Print("🗑️ Nom de la branche à supprimer: ")
	branchName := ga.getUserInput()
	if branchName == "" {
		return fmt.Errorf("nom requis")
	}
	
	current := ga.getCurrentBranch()
	if branchName == current {
		fmt.Println("❌ Impossible de supprimer la branche courante")
		return nil
	}
	
	fmt.Printf("⚠️ Êtes-vous sûr de vouloir supprimer '%s'? (o/N): ", branchName)
	if strings.ToLower(ga.getUserInput()) != "o" {
		fmt.Println("❌ Suppression annulée")
		return nil
	}
	
	_, err := ga.runCommand("git", "branch", "-d", branchName)
	if err != nil {
		// Essayer force delete
		fmt.Print("⚠️ Branche non fusionnée. Forcer la suppression? (o/N): ")
		if strings.ToLower(ga.getUserInput()) == "o" {
			_, err = ga.runCommand("git", "branch", "-D", branchName)
		}
	}
	
	if err != nil {
		return err
	}
	
	fmt.Printf("✅ Branche '%s' supprimée!\n", branchName)
	ga.addToHistory(fmt.Sprintf("Branche supprimée: %s", branchName))
	return nil
}

func (ga *GitAssistant) mergeBranch() error {
	current := ga.getCurrentBranch()
	fmt.Printf("🔀 Fusion vers la branche courante (%s)\n", current)
	fmt.Print("Nom de la branche à fusionner: ")
	
	branchName := ga.getUserInput()
	if branchName == "" {
		return fmt.Errorf("nom requis")
	}
	
	_, err := ga.runCommand("git", "merge", branchName)
	if err != nil {
		fmt.Println("❌ Conflit détecté! Résolvez manuellement puis recommitez.")
		return err
	}
	
	fmt.Printf("✅ Branche '%s' fusionnée dans '%s'!\n", branchName, current)
	ga.addToHistory(fmt.Sprintf("Fusion: %s → %s", branchName, current))
	return nil
}

func (ga *GitAssistant) showHistory() error {
	if len(ga.lastActions) == 0 {
		fmt.Println("📜 Aucune action récente")
		return nil
	}
	
	fmt.Println("📜 === HISTORIQUE DES ACTIONS ===")
	for i := len(ga.lastActions) - 1; i >= 0; i-- {
		fmt.Printf("%d. %s\n", len(ga.lastActions)-i, ga.lastActions[i])
	}
	
	return nil
}

func (ga *GitAssistant) interactiveLog() error {
	fmt.Printf("📜 === %s ===\n", bold("HISTORIQUE INTERACTIF"))
	
	output, err := ga.runCommand("git", "log", "--oneline", "-15", "--graph", "--decorate")
	if err != nil {
		return err
	}
	
	// Améliorer l'affichage avec des séparations après chaque commit
	lines := strings.Split(strings.TrimSpace(output), "\n")
	fmt.Println()
	for i, line := range lines {
		fmt.Printf("  %s\n", line)
		// Ajouter une ligne de séparation après chaque commit (sauf le dernier)
		if i < len(lines)-1 {
			fmt.Printf("  %s\n", cyan(strings.Repeat("─", 80)))
		}
	}
	fmt.Println()
	
	fmt.Printf("%s:\n", cyan("Actions disponibles"))
	fmt.Println("1. 👀 Voir détails d'un commit")
	fmt.Println("2. ⏪ Reset vers un commit")
	fmt.Println("3. 🌱 Créer branche depuis commit")
	fmt.Println("4. 🔍 Rechercher dans l'historique")
	fmt.Print(cyan("\nChoisissez (1-4): "))
	
	choice := ga.getUserInput()
	
	switch choice {
	case "1":
		return ga.showCommitDetails()
	case "2":
		return ga.resetToCommit()
	case "3":
		return ga.createBranchFromCommit()
	case "4":
		return ga.searchInHistory()
	}
	
	return nil
}

func (ga *GitAssistant) showCommitDetails() error {
	fmt.Print(cyan("🔍 Hash du commit: "))
	hash := ga.getUserInput()
	
	if hash == "" {
		return fmt.Errorf("hash requis")
	}
	
	// Afficher les informations générales du commit
	output, err := ga.runCommand("git", "show", "--stat", "--pretty=format:%h - %s%n%an <%ae>%n%ad%n", hash)
	if err != nil {
		return err
	}
	
	fmt.Printf("📋 %s:\n", cyan("Détails du commit"))
	fmt.Println(output)
	
	// Afficher directement le diff complet
	diffOutput, err := ga.runCommand("git", "diff", hash+"^", hash)
	if err != nil {
		// Si pas de parent (premier commit), utiliser git show
		diffOutput, err = ga.runCommand("git", "show", "--format=", hash)
		if err != nil {
			return fmt.Errorf("impossible d'obtenir le diff: %v", err)
		}
	}
	
	if strings.TrimSpace(diffOutput) == "" {
		fmt.Println("Aucun changement de fichier dans ce commit")
		return nil
	}
	
	fmt.Printf("\n🔍 %s:\n", cyan("Diff complet"))
	ga.displayColoredDiff(diffOutput)
	
	return nil
}

func (ga *GitAssistant) displayColoredDiff(diffText string) {
	lines := strings.Split(diffText, "\n")
	
	for _, line := range lines {
		// Colorier selon le préfixe de la ligne
		switch {
		case strings.HasPrefix(line, "diff --git"):
			// Headers de diff
			fmt.Println(bold(line))
		case strings.HasPrefix(line, "index "):
			// Index des fichiers
			fmt.Println(cyan(line))
		case strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---"):
			// Headers de fichiers
			fmt.Println(cyan(line))
		case strings.HasPrefix(line, "@@"):
			// Numéros de lignes
			fmt.Println(cyan(line))
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			// Lignes ajoutées
			fmt.Println(green(line))
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			// Lignes supprimées
			fmt.Println(red(line))
		case line == "":
			// Lignes vides
			fmt.Println()
		default:
			// Lignes de contexte (inchangées) - commencent par un espace
			fmt.Println(line)
		}
	}
}

func (ga *GitAssistant) searchInHistory() error {
	fmt.Print("🔍 Rechercher (message/auteur/fichier): ")
	query := ga.getUserInput()
	
	if query == "" {
		return fmt.Errorf("terme de recherche requis")
	}
	
	// Recherche dans les messages
	output, _ := ga.runCommand("git", "log", "--oneline", "--grep="+query, "-i")
	if output != "" {
		fmt.Println("📝 Commits avec ce message:")
		fmt.Println(output)
	}
	
	// Recherche par fichier
	output2, _ := ga.runCommand("git", "log", "--oneline", "--", "*"+query+"*")
	if output2 != "" {
		fmt.Println("📁 Commits affectant ce fichier:")
		fmt.Println(output2)
	}
	
	if output == "" && output2 == "" {
		fmt.Println("❌ Aucun résultat trouvé")
	}
	
	return nil
}

func (ga *GitAssistant) projectInsights() error {
	fmt.Printf("📊 === %s ===\n", bold("ANALYSE DU PROJET"))
	
	// Statistiques générales
	commits, files, branches := ga.getRepoStats()
	fmt.Printf("📈 Statistiques:\n")
	fmt.Printf("  • %s commits au total\n", green(strconv.Itoa(commits)))
	fmt.Printf("  • %s fichiers suivis\n", green(strconv.Itoa(files)))
	fmt.Printf("  • %s branches\n\n", green(strconv.Itoa(branches)))
	
	// Liste des branches avec détails
	fmt.Printf("🌿 %s:\n", cyan("Branches disponibles"))
	branchesOutput, err := ga.runCommand("git", "branch", "-v")
	if err == nil && branchesOutput != "" {
		lines := strings.Split(strings.TrimSpace(branchesOutput), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				if strings.HasPrefix(line, "*") {
					fmt.Printf("  → %s %s\n", green(line[2:]), cyan("(branche actuelle)"))
				} else {
					fmt.Printf("  • %s\n", line)
				}
			}
		}
	} else {
		fmt.Println("  Aucune branche trouvée")
	}
	fmt.Println()
	
	// Analyse des extensions
	output, err := ga.runCommand("git", "ls-files")
	if err == nil && output != "" {
		ga.analyzeFileTypes(strings.Split(output, "\n"))
	}
	
	// Activité récente
	recentOutput, err := ga.runCommand("git", "log", "--since=1.week.ago", "--oneline")
	if err == nil {
		recentCommits := len(strings.Split(strings.TrimSpace(recentOutput), "\n"))
		if recentOutput == "" {
			recentCommits = 0
		}
		fmt.Printf("⚡ Activité récente: %s commits cette semaine\n", green(strconv.Itoa(recentCommits)))
	}
	
	// Taille du dépôt
	sizeOutput, _ := ga.runCommand("git", "count-objects", "-vH")
	if sizeOutput != "" {
		lines := strings.Split(sizeOutput, "\n")
		for _, line := range lines {
			if strings.Contains(line, "size-pack") {
				fmt.Printf("💾 Taille: %s\n", green(strings.Split(line, " ")[1]))
				break
			}
		}
	}
	
	return nil
}

func (ga *GitAssistant) analyzeFileTypes(files []string) {
	extCount := make(map[string]int)
	
	for _, file := range files {
		if file == "" {
			continue
		}
		ext := filepath.Ext(file)
		if ext == "" {
			ext = "sans extension"
		}
		extCount[ext]++
	}
	
	// Trier par nombre
	type kv struct {
		Key   string
		Value int
	}
	
	var sorted []kv
	for k, v := range extCount {
		sorted = append(sorted, kv{k, v})
	}
	
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	
	fmt.Println("📂 Types de fichiers:")
	for i, kv := range sorted {
		if i >= 5 {
			break
		}
		fmt.Printf("  • %s: %d fichiers\n", kv.Key, kv.Value)
	}
	fmt.Println()
}

// Fonctions existantes simplifiées
func (ga *GitAssistant) initRepo() error {
	fmt.Println("🔧 Initialisation du dépôt Git...")
	_, err := ga.runCommand("git", "init")
	if err != nil {
		return fmt.Errorf("erreur lors de l'initialisation: %v", err)
	}
	fmt.Println("✅ Dépôt Git initialisé avec succès!")
	ga.addToHistory("Dépôt initialisé")
	return nil
}

func (ga *GitAssistant) getStatus() (string, error) {
	return ga.runCommand("git", "status", "--porcelain")
}

func (ga *GitAssistant) addAll() error {
	fmt.Println("📝 Ajout de tous les fichiers...")
	_, err := ga.runCommand("git", "add", ".")
	if err != nil {
		return fmt.Errorf("erreur lors de l'ajout des fichiers: %v", err)
	}
	fmt.Println("✅ Fichiers ajoutés!")
	return nil
}

func (ga *GitAssistant) commit(message string) error {
	if message == "" {
		message = fmt.Sprintf("Auto-commit: %s", time.Now().Format("2006-01-02 15:04:05"))
	}
	
	fmt.Printf("💾 Commit avec le message: %s\n", message)
	_, err := ga.runCommand("git", "commit", "-m", message)
	if err != nil {
		return fmt.Errorf("erreur lors du commit: %v", err)
	}
	fmt.Println("✅ Commit effectué!")
	return nil
}

func (ga *GitAssistant) autoSync() error {
	fmt.Println("🔄 Synchronisation automatique...")
	
	status, err := ga.getStatus()
	if err != nil {
		return err
	}
	
	if strings.TrimSpace(status) == "" {
		fmt.Println("ℹ️ Aucun changement détecté")
		return nil
	}
	
	if err := ga.addAll(); err != nil {
		return err
	}
	
	if err := ga.commit(""); err != nil {
		return err
	}
	
	fmt.Println("🎉 Synchronisation locale terminée!")
	ga.addToHistory("Synchronisation automatique")
	return nil
}

func (ga *GitAssistant) switchBranch(name string) error {
	if name == "" {
		fmt.Print("🔄 Nom de la branche: ")
		name = ga.getUserInput()
	}
	
	if name == "" {
		return fmt.Errorf("nom de branche requis")
	}
	
	fmt.Printf("🔄 Changement vers la branche: %s\n", name)
	_, err := ga.runCommand("git", "checkout", name)
	if err != nil {
		return fmt.Errorf("erreur lors du changement de branche: %v", err)
	}
	fmt.Println("✅ Branche changée!")
	ga.addToHistory(fmt.Sprintf("Changement vers: %s", name))
	return nil
}

func (ga *GitAssistant) resetToCommit() error {
	fmt.Println("📜 Historique récent:")
	output, err := ga.runCommand("git", "log", "--oneline", "-10")
	if err != nil {
		return err
	}
	fmt.Println(output)
	
	fmt.Println("\n🔄 Types de reset:")
	fmt.Println("1. 🟢 SOFT - Garde les changements dans le staging")
	fmt.Println("2. 🟡 MIXED - Garde les changements mais pas dans le staging")
	fmt.Println("3. 🔴 HARD - Supprime TOUS les changements")
	fmt.Print("\nType (1-3): ")
	
	resetType := ga.getUserInput()
	var resetFlag string
	
	switch resetType {
	case "1":
		resetFlag = "--soft"
	case "2":
		resetFlag = "--mixed"
	case "3":
		resetFlag = "--hard"
	default:
		resetFlag = "--mixed"
	}
	
	fmt.Print("🎯 Hash du commit (ou HEAD~n): ")
	commitHash := ga.getUserInput()
	
	if commitHash == "" {
		return fmt.Errorf("hash requis")
	}
	
	_, err = ga.runCommand("git", "reset", resetFlag, commitHash)
	if err != nil {
		return err
	}
	
	fmt.Println("✅ Reset effectué!")
	ga.addToHistory(fmt.Sprintf("Reset %s vers %s", resetFlag, commitHash))
	return nil
}

func (ga *GitAssistant) createBranchFromCommit() error {
	fmt.Println("📜 Historique récent:")
	output, err := ga.runCommand("git", "log", "--oneline", "-10")
	if err != nil {
		return err
	}
	fmt.Println(output)
	
	fmt.Print("\n🎯 Hash du commit: ")
	commitHash := ga.getUserInput()
	
	fmt.Print("🌱 Nom de la branche: ")
	branchName := ga.getUserInput()
	
	if commitHash == "" || branchName == "" {
		return fmt.Errorf("hash et nom requis")
	}
	
	_, err = ga.runCommand("git", "checkout", "-b", branchName, commitHash)
	if err != nil {
		return err
	}
	
	fmt.Println("✅ Branche créée et activée!")
	ga.addToHistory(fmt.Sprintf("Branche %s depuis %s", branchName, commitHash))
	return nil
}

func (ga *GitAssistant) setWorkingDirectory(newPath string) error {
	if newPath == "" {
		return fmt.Errorf("chemin vide")
	}
	
	// Convertir en chemin absolu
	absPath, err := filepath.Abs(newPath)
	if err != nil {
		return fmt.Errorf("chemin invalide: %v", err)
	}
	
	// Vérifier que le répertoire existe
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("le répertoire n'existe pas: %s", absPath)
	}
	
	ga.workingDir = absPath
	ga.lastActions = make([]string, 0) // Reset l'historique pour le nouveau projet
	fmt.Printf(green("✅ Répertoire défini: %s\n"), ga.workingDir)
	return nil
}

func (ga *GitAssistant) changeDirectory() error {
	fmt.Print("📁 Nouveau répertoire: ")
	newPath := ga.getUserInput()
	return ga.setWorkingDirectory(newPath)
}

func (ga *GitAssistant) handleNonGitRepo() bool {
	fmt.Println("\n⚠️ Ce répertoire n'est pas un dépôt Git.")
	fmt.Println("1. 🔧 Initialiser un dépôt Git ici")
	fmt.Println("2. 📁 Changer de répertoire")
	fmt.Println("3. ❌ Quitter")
	fmt.Print("\nChoisissez (1-3): ")
	
	choice := ga.getUserInput()
	
	switch choice {
	case "1":
		if err := ga.initRepo(); err != nil {
			fmt.Printf("❌ Erreur: %v\n", err)
			return false
		}
		return true
	case "2":
		if err := ga.changeDirectory(); err != nil {
			fmt.Printf("❌ Erreur: %v\n", err)
			return false
		}
		return ga.checkGitRepoOrHandle()
	case "3":
		fmt.Println("👋 Au revoir!")
		return false
	default:
		fmt.Println("❌ Choix invalide!")
		return false
	}
}

func (ga *GitAssistant) checkGitRepoOrHandle() bool {
	if ga.isGitRepo() {
		return true
	}
	return ga.handleNonGitRepo()
}

func (ga *GitAssistant) clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (ga *GitAssistant) showSmartMenu() {
	ga.clearScreen()
	
	// ASCII Art titre
	fmt.Println(cyan(`
   ___ ___ _____ ___ _____ ___ _    
  / __|_ _|_   _/ __|_   _| _ \ |   
 | (_ || |  | || (__  | | |   / |__ 
  \___|___| |_| \___| |_| |_|_\____|`))
	
	// Header avec infos contextuelles
	fmt.Printf("\n🚀 === %s ===\n", bold("GIT ASSISTANT INTELLIGENT"))
	fmt.Printf("📁 Répertoire: %s\n", cyan(ga.workingDir))
	
	if ga.isGitRepo() {
		branch := ga.getCurrentBranch()
		commits, _, _ := ga.getRepoStats()
		fmt.Printf("🌿 Branche: %s | 📊 %d commits", green(branch), commits)
		
		// Vérifier s'il y a des changements
		status, _ := ga.getStatus()
		if strings.TrimSpace(status) != "" {
			fmt.Print(" | " + red("⚠️ Changements non commitées"))
		}
		fmt.Println()
		
		fmt.Printf("\n=== %s ===\n", cyan("ACTIONS RAPIDES"))
		fmt.Println("1. ⚡ Commit rapide (messages prédéfinis)")
		
		fmt.Printf("\n=== %s ===\n", cyan("GESTION AVANCÉE"))
		fmt.Println("2. 🌿 Gestion intelligente des branches")
		fmt.Println("3. 📜 Historique interactif")
		fmt.Println("4. 📊 Analyse du projet")
		
		fmt.Printf("\n=== %s ===\n", cyan("NAVIGATION"))
		fmt.Println("5. 📁 Changer de répertoire")
		fmt.Println("6. 🔧 Initialiser Git")
	} else {
		fmt.Println(red("⚠️ Pas un dépôt Git"))
		fmt.Printf("\n=== %s ===\n", cyan("ACTIONS DISPONIBLES"))
		fmt.Println("1. 🔧 Initialiser un dépôt Git ici")
		fmt.Println("2. 📁 Changer de répertoire")
	}
	
	fmt.Println("\n0. ❌ Quitter")
	fmt.Print(cyan("\n💫 Choisissez une action: "))
}

func (ga *GitAssistant) getUserInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func (ga *GitAssistant) run() {
	fmt.Println(bold("🎯 Git Assistant Intelligent démarré!"))
	
	// Première action obligatoire : définir le répertoire de travail
	fmt.Printf("\n📁 === %s ===\n", cyan("SÉLECTION DU RÉPERTOIRE DE TRAVAIL"))
	fmt.Printf("Répertoire actuel: %s\n", cyan(ga.workingDir))
	fmt.Print(cyan("Entrez le chemin du dossier de travail: "))
	
	newPath := ga.getUserInput()
	if newPath != "" {
		if err := ga.setWorkingDirectory(newPath); err != nil {
			fmt.Printf(red("❌ Erreur: %v\n"), err)
			fmt.Print("Continuer avec le répertoire actuel? (o/N): ")
			if strings.ToLower(ga.getUserInput()) != "o" {
				fmt.Println("👋 Au revoir!")
				return
			}
		}
	}
	
	for {
		ga.showSmartMenu()
		choice := ga.getUserInput()
		
		switch choice {
		case "1":
			if ga.isGitRepo() {
				if err := ga.quickCommit(); err != nil {
					fmt.Printf(red("❌ Erreur: %v\n"), err)
				}
			} else {
				if err := ga.initRepo(); err != nil {
					fmt.Printf(red("❌ Erreur: %v\n"), err)
				}
			}
			
		case "2":
			if ga.isGitRepo() {
				if err := ga.intelligentBranching(); err != nil {
					fmt.Printf(red("❌ Erreur: %v\n"), err)
				}
			} else {
				if err := ga.changeDirectory(); err != nil {
					fmt.Printf(red("❌ Erreur: %v\n"), err)
				}
			}
			
		case "3":
			if ga.isGitRepo() {
				if err := ga.interactiveLog(); err != nil {
					fmt.Printf(red("❌ Erreur: %v\n"), err)
				}
			} else {
				fmt.Println(red("❌ Cette action nécessite un dépôt Git"))
			}
			
		case "4":
			if ga.isGitRepo() {
				if err := ga.projectInsights(); err != nil {
					fmt.Printf(red("❌ Erreur: %v\n"), err)
				}
			} else {
				fmt.Println(red("❌ Cette action nécessite un dépôt Git"))
			}
			
		case "5":
			if ga.isGitRepo() {
				if err := ga.changeDirectory(); err != nil {
					fmt.Printf(red("❌ Erreur: %v\n"), err)
				}
			} else {
				fmt.Println(red("❌ Cette action nécessite un dépôt Git"))
			}
			
		case "6":
			if ga.isGitRepo() {
				if err := ga.initRepo(); err != nil {
					fmt.Printf(red("❌ Erreur: %v\n"), err)
				}
			} else {
				fmt.Println(red("❌ Cette action nécessite un dépôt Git"))
			}
			
		case "0":
			fmt.Println("👋 Au revoir!")
			return
			
		default:
			fmt.Println(red("❌ Option invalide!"))
		}
		
		fmt.Println("\n⏸️ Appuyez sur Entrée pour continuer...")
		ga.getUserInput()
	}
}

func main() {
	assistant := NewGitAssistant()
	assistant.run()
}
		