const STEPS = [
  {
    step: 1,
    title: 'Explora',
    description: 'Navega por los eventos disponibles y encuentra el que más te guste.',
  },
  {
    step: 2,
    title: 'Compra',
    description: 'Selecciona tus entradas y completa la compra de forma segura.',
  },
  {
    step: 3,
    title: 'Disfruta',
    description: 'Recibe tus entradas digitales y disfruta del evento.',
  },
]

export default function HowItWorks() {
  return (
    <section className="how-it-works-section">
      <div className="section-header">
        <h2>¿Cómo funciona?</h2>
      </div>
      <div className="steps-grid">
        {STEPS.map((step) => (
          <div key={step.step} className="step-card">
            <div className="step-number">{step.step}</div>
            <h3>{step.title}</h3>
            <p>{step.description}</p>
          </div>
        ))}
      </div>
    </section>
  )
}
