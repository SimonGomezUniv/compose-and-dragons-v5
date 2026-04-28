const fs = require('fs').promises;
const fsSync = require('fs');
const path = require('path');
const yaml = require('yaml');

/**
 * Charger tous les tickets YAML depuis un répertoire
 * @param {string} ticketsDir - Chemin du répertoire
 * @returns {Promise<Array>} Array de tickets
 */
async function loadTicketsFromYAML(ticketsDir) {
  const tickets = [];

  try {
    if (!fsSync.existsSync(ticketsDir)) {
      console.warn(`Tickets directory not found: ${ticketsDir}`);
      return tickets;
    }

    const files = await fs.readdir(ticketsDir);
    const ymlFiles = files.filter((f) => f.endsWith('.yml') || f.endsWith('.yaml'));

    for (const file of ymlFiles) {
      try {
        const filePath = path.join(ticketsDir, file);
        const content = await fs.readFile(filePath, 'utf-8');
        const ticket = yaml.parse(content);

        // Valider champs requis
        if (ticket.id && ticket.title && ticket.description) {
          tickets.push({
            ...ticket,
            createdAt: new Date().toISOString(),
          });
          console.log(`✅ Loaded ticket: ${ticket.id}`);
        } else {
          console.warn(`⚠️  Invalid ticket in ${file}: missing required fields`);
        }
      } catch (error) {
        console.error(`Error loading ${file}:`, error.message);
      }
    }
  } catch (error) {
    console.error('Error reading tickets directory:', error.message);
  }

  return tickets;
}

/**
 * Charger contexte d'expertise depuis fichier MD
 * @param {string} contextPath - Chemin du fichier MD
 * @returns {Promise<string>} Contenu du fichier
 */
async function loadExpertContext(contextPath) {
  try {
    if (!fsSync.existsSync(contextPath)) {
      console.warn(`Expert context file not found: ${contextPath}`);
      return '';
    }

    const content = await fs.readFile(contextPath, 'utf-8');
    return content;
  } catch (error) {
    console.error('Error loading expert context:', error.message);
    return '';
  }
}

module.exports = { loadTicketsFromYAML, loadExpertContext };
