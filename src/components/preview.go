package components

import (
	"html"

	"github.com/RyannKim327/constellation/config/utils"
)

// PreviewPanel renders the live Markdown preview target for a textarea id.
func PreviewPanel(inputID string) string {
	return utils.HTParser(`
		<div class="preview-panel">
			<span class="eyebrow">Live preview</span>
			<div class="preview-body" data-preview-for="{input}" aria-live="polite"></div>
		</div>
	`, utils.PARAMS{"input": html.EscapeString(inputID)})
}

// LivePreviewScript wires every preview panel on the page to the Markdown API.
func LivePreviewScript() string {
	return `
		<script>
		(function () {
			var timers = {};
			var panels = document.querySelectorAll('[data-preview-for]');
			panels.forEach(function (panel) {
				var input = document.getElementById(panel.getAttribute('data-preview-for'));
				if (!input) return;
				function render() {
					fetch('/api/markdown', {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify({ text: input.value })
					}).then(function (response) { return response.json(); })
						.then(function (data) { panel.innerHTML = data.html || ''; })
						.catch(function () { panel.textContent = 'Preview unavailable.'; });
				}
				input.addEventListener('input', function () {
					clearTimeout(timers[input.id]);
					timers[input.id] = setTimeout(render, 150);
				});
				render();
			});
		})();
		</script>
	`
}
