(function() {
  var els = [];
  els.push.apply(els, document.getElementsByTagName("h2"));
  els.push.apply(els, document.getElementsByTagName("h3"));
  els.push.apply(els, document.getElementsByTagName("h4"));
  els.push.apply(els, document.getElementsByTagName("h5"));
  for (var i=0; i<els.length; i++) {
    (function(el) {
      if (el.id) {
        var node;
        el.addEventListener("mouseenter", function(evt) {
          if (!node) {
            node = document.createElement("a");
            node.className = "section-link";
            node.href = "#" + el.id;
            node.textContent = "§";
            el.appendChild(node);
          }
        });
        el.addEventListener("mouseleave", function(evt) {
          if (node) {
            el.removeChild(node);
            node = null;
          }
        });
      }
    })(els[i]);
  }

})();
