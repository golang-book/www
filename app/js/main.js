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


var mss = document.getElementsByClassName("multi-step");
for (var i=0; i<mss.length; i++) {
  var ms = mss[i];

  var h3s = ms.getElementsByTagName("h3");
  var bc = document.createElement("ul");
  bc.className = "step-breadcrumbs";
  ms.insertBefore(bc, ms.children[1]);

  for (var j=0; j<h3s.length; j++) {
    var h3 = h3s[j];

    // add the breadcrumb
    var li = document.createElement("li");
    var a = document.createElement("a");
    a.setAttribute("href", "#" + h3.id);
    a.textContent = h3.textContent;
    li.appendChild(a);
    bc.appendChild(li);
  }

  for (var j=1; j<=h3s.length; j++) {
    var footer = document.createElement("div");
    footer.className = "step-footer";

    if (j > 1) {
      var a = document.createElement("a");
      a.className = "previous";
      a.setAttribute("href", "#" + h3s[j-2].id);
      a.innerHTML = "&larr; Previous";
      footer.appendChild(a);
    }

    if (j > 0 && j < h3s.length) {
      // add a next link before the heading
      var a = document.createElement("a");
      a.className = "next";
      a.setAttribute("href", "#" + h3s[j].id);
      a.innerHTML = "Next &rarr;";
      footer.appendChild(a);
    }

    if (j < h3s.length) {
      h3s[j].parentNode.insertBefore(footer, h3s[j]);
    } else {
      ms.appendChild(footer);
    }
  }
}

function updateMultiStep() {
  for (var i=0; i<mss.length; i++) {
    var ms = mss[i];
    if (ms.getAttribute("data-for")) {
      ms.style.display = location.hash.indexOf(ms.getAttribute("data-for")) === 1
        ? ""
        : "none";
    }

    // find the hash (either user-entered or the first step)
    var hash = null;
    for (var j=0; j<ms.children.length; j++) {
      var el = ms.children[j];
      if (el.tagName === "H3") {
        if (location.hash.substr(1) === el.id) {
          hash = el.id;
        } else if (!hash) {
          hash = el.id;
        }
      }
    }

    // update the breadcrumbs
    var bc = ms.getElementsByClassName("step-breadcrumbs")[0];
    if (bc) {
      var lis = bc.getElementsByTagName("li");
      for (var j=0; j<lis.length; j++) {
        var li = lis[j];
        if (li.firstChild.getAttribute("href") === "#" + hash) {
          li.className = "selected";
        } else {
          li.className = "";
        }
      }
    }

    var show = true;
    for (var j=0; j<ms.children.length; j++) {
      var el = ms.children[j];
      if (el.tagName === "H3") {
        show = el.id === hash;
        console.log("EL", el.id===hash);
      }
      if (show) {
        el.style.display = "";
      } else{
        el.style.display = "none";
      }
    }
  }
}
updateMultiStep();
window.addEventListener("hashchange", function() {
  updateMultiStep();
}, false);

})();
