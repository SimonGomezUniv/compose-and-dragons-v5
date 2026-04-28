const axios = require('axios');

async function testEmbedding() {
  try {
    const response = await axios.post('http://localhost:11435/v1/embeddings', {
      model: 'embeddinggemma:latest',
      input: 'test embedding',
    });
    
    console.log('Response structure:');
    console.log(JSON.stringify(response.data, null, 2));
  } catch (error) {
    console.error('Error:', error.message);
    if (error.response) {
      console.error('Response:', error.response.data);
    }
  }
}

testEmbedding();
