// Landing page JavaScript
document.addEventListener('DOMContentLoaded', function() {
    // Add some interactive effects
    const coins = document.querySelectorAll('.coin');

    coins.forEach((coin, index) => {
        coin.addEventListener('mouseenter', function() {
            this.style.transform = 'scale(1.2) rotate(360deg)';
            this.style.transition = 'transform 0.5s ease';
        });

        coin.addEventListener('mouseleave', function() {
            this.style.transform = 'scale(1) rotate(0deg)';
        });
    });

    // Smooth scroll for CTA button
    const ctaButton = document.querySelector('.cta-button');
    if (ctaButton) {
        ctaButton.addEventListener('click', function(e) {
            // Add a loading effect
            this.textContent = 'Loading Dashboard...';
            this.style.opacity = '0.7';
        });
    }

    // Add some particle effect or animation
    function createParticle() {
        const particle = document.createElement('div');
        particle.className = 'particle';
        particle.style.cssText = `
            position: absolute;
            width: 4px;
            height: 4px;
            background: rgba(255, 255, 255, 0.5);
            border-radius: 50%;
            pointer-events: none;
            animation: particleFloat 10s linear infinite;
            left: ${Math.random() * 100}%;
            top: 100vh;
        `;

        document.body.appendChild(particle);

        setTimeout(() => {
            particle.remove();
        }, 10000);
    }

    // Create particles every few seconds
    setInterval(createParticle, 2000);

    // Add CSS for particle animation
    const style = document.createElement('style');
    style.textContent = `
        @keyframes particleFloat {
            0% { transform: translateY(0px) rotate(0deg); opacity: 1; }
            100% { transform: translateY(-100vh) rotate(360deg); opacity: 0; }
        }
    `;
    document.head.appendChild(style);
});
