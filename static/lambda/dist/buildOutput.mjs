export const handler = async (event) => {
    const response = {
      source: {
        bucket: event.detail.bucket.name,
        file: event.detail.object.key
      },
      labels: event.Rekognition.Labels,
      summary: event.Bedrock.Body.results[0].outputText
    }
    
    return response;
  };