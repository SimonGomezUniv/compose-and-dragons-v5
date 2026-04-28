/**
 * @typedef {Object} SupportTicket
 * @property {string} id - Identifiant unique
 * @property {string} title - Titre du ticket
 * @property {string} description - Description détaillée
 * @property {string} [category] - Catégorie (optionnelle)
 * @property {'low'|'medium'|'high'|'critical'} [priority] - Priorité
 * @property {string[]} [tags] - Tags pour recherche
 * @property {Date} [created_at] - Date de soumission
 * @property {string} [resolution] - Résolution si trouvée
 */

/**
 * @typedef {Object} AnalysisResult
 * @property {string} ticketId
 * @property {string} suggestedCategory
 * @property {number} confidence - Entre 0 et 1
 * @property {string[]} suggestions - Suggestions de résolution
 * @property {string[]} additionalInfoNeeded - Infos manquantes
 * @property {string} reasoning - Explication brève
 * @property {RelatedTicket[]} relatedPastTickets - Tickets similaires
 */

/**
 * @typedef {Object} RelatedTicket
 * @property {string} id
 * @property {string} title
 * @property {number} similarity - Score de similarité
 * @property {string} [resolution]
 */

/**
 * Valider un ticket
 * @param {SupportTicket} ticket
 * @returns {boolean}
 * @throws {Error} Si champs requis manquent
 */
function validateTicket(ticket) {
  if (!ticket.id || !ticket.title || !ticket.description) {
    throw new Error('Missing required fields: id, title, description');
  }
  return true;
}

/**
 * Nettoyer description pour embedding
 * @param {string} text
 * @returns {string}
 */
function cleanText(text) {
  return text
    .toLowerCase()
    .replace(/[^\w\s]/g, ' ')
    .replace(/\s+/g, ' ')
    .trim();
}

module.exports = { validateTicket, cleanText };
