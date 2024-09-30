export const handler = async (event) => {
    const confidenceLevel = parseInt(process.env.CONFIDENCE_LEVEL) || 90;
  
    const labels = event.Rekognition.Labels;
  
    const filteredLabels = labels
      .filter((label) => label.Confidence > confidenceLevel)
      .map((label) =>
        label.Instances.length > 0
          ? `${label.Instances.length} ${label.Name}`
          : label.Name
      )
      .join(", ");
  
    const response = {
      labels: filteredLabels,
    };
    return response;
  };