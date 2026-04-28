/**
 * Calculer la similarité cosinus entre deux vecteurs
 * @param {number[]} a
 * @param {number[]} b
 * @returns {number} Score entre 0 et 1
 */
function cosineSimilarity(a, b) {
  if (!a || !b || a.length !== b.length) {
    throw new Error('Vectors must have the same length');
  }

  let dotProduct = 0;
  let magnitudeA = 0;
  let magnitudeB = 0;

  for (let i = 0; i < a.length; i++) {
    dotProduct += a[i] * b[i];
    magnitudeA += a[i] * a[i];
    magnitudeB += b[i] * b[i];
  }

  magnitudeA = Math.sqrt(magnitudeA);
  magnitudeB = Math.sqrt(magnitudeB);

  if (magnitudeA === 0 || magnitudeB === 0) {
    return 0;
  }

  return dotProduct / (magnitudeA * magnitudeB);
}

/**
 * Rechercher les records similaires
 * @param {Array} records - Tous les records du vector store
 * @param {number[]} queryEmbedding - Embedding de la requête
 * @param {number} threshold - Seuil de similarité (défaut: 0.4)
 * @param {number} topN - Nombre de résultats (défaut: 5)
 * @returns {Array} Records triés par similarité descendante
 */
function searchSimilar(records, queryEmbedding, threshold = 0.4, topN = 5) {
  if (!records || records.length === 0) {
    return [];
  }

  // Calculer similarité pour chaque record
  const scored = records.map((record) => ({
    ...record,
    similarity: cosineSimilarity(queryEmbedding, record.embedding),
  }));

  // Filtrer par threshold et trier
  return scored
    .filter((record) => record.similarity >= threshold)
    .sort((a, b) => b.similarity - a.similarity)
    .slice(0, topN);
}

module.exports = { cosineSimilarity, searchSimilar };
