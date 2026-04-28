const { v4: uuidv4 } = require('uuid');

/**
 * @typedef {Object} VectorRecord
 * @property {string} id - UUID unique du record
 * @property {string} ticketId - ID du ticket
 * @property {string} content - Contenu enrichi avec metadata
 * @property {number[]} embedding - Vecteur d'embedding (1024 dimensions)
 * @property {Object} metadata - Métadonnées extraites
 * @property {string[]} metadata.keywords - 4 mots-clés
 * @property {string} metadata.mainCategory - Catégorie principale
 * @property {string} metadata.subcategory - Sous-catégorie
 * @property {'low'|'medium'|'high'} metadata.importance - Importance
 * @property {Date} timestamp - Timestamp création
 */

class VectorStore {
  constructor() {
    this.records = new Map();
  }

  /**
   * Charger le store depuis un fichier JSON
   * @param {string} path - Chemin du fichier
   */
  async loadFromFile(path) {
    const fs = require('fs').promises;
    try {
      const fsSync = require('fs');
      if (fsSync.existsSync(path)) {
        const data = JSON.parse(await fs.readFile(path, 'utf-8'));
        this.records = new Map(Object.entries(data));
        console.log(`✅ Loaded ${this.records.size} vector records from ${path}`);
      }
    } catch (error) {
      console.error('Error loading vector store:', error.message);
    }
  }

  /**
   * Sauvegarder le store en JSON
   * @param {string} path
   */
  async saveToFile(path) {
    const fs = require('fs').promises;
    const data = Object.fromEntries(this.records);
    await fs.writeFile(path, JSON.stringify(data, null, 2));
    console.log(`✅ Saved ${this.records.size} vector records to ${path}`);
  }

  /**
   * Ajouter un record
   * @param {Omit<VectorRecord, 'id'>} record
   * @returns {VectorRecord}
   */
  addRecord(record) {
    const fullRecord = {
      id: uuidv4(),
      ...record,
      timestamp: new Date().toISOString(),
    };
    this.records.set(fullRecord.id, fullRecord);
    return fullRecord;
  }

  /**
   * Obtenir tous les records
   * @returns {VectorRecord[]}
   */
  getAllRecords() {
    return Array.from(this.records.values());
  }

  /**
   * Obtenir un record par ID
   * @param {string} id
   * @returns {VectorRecord|undefined}
   */
  getRecord(id) {
    return this.records.get(id);
  }

  /**
   * Supprimer un record
   * @param {string} id
   */
  deleteRecord(id) {
    return this.records.delete(id);
  }

  /**
   * Nombre de records stockés
   * @returns {number}
   */
  size() {
    return this.records.size;
  }
}

module.exports = { VectorStore };
